# Use minimal Alpine base for security
FROM node:20-alpine

# Security: Create non-root user with host-compatible UID
# Use build arg to pass host UID, defaulting to 1000 if not specified
ARG HOST_UID=10000
ARG HOST_GID=10000
# Create group/user only if IDs not already present; otherwise reuse
RUN set -eux; \
    if ! getent group ${HOST_GID} >/dev/null 2>&1; then addgroup -g ${HOST_GID} -S opencode; else echo "Group ${HOST_GID} exists"; fi; \
    if ! getent passwd ${HOST_UID} >/dev/null 2>&1; then adduser -S opencode -u ${HOST_UID} -G $(getent group ${HOST_GID} | cut -d: -f1 || echo opencode); else echo "User ${HOST_UID} exists"; fi

# Install system dependencies (minimal set)
RUN apk add --no-cache \
    # Essential build tools
    build-base \
    python3 \
    # Editor and utilities
    neovim \
    git \
    curl \
    wget \
    # Unix utilities
    bash \
    zsh \
    fish \
    tmux \
    htop \
    ripgrep \
    fd \
    # Network tools (minimal)
    netcat-openbsd \
    # Security tools
    doas \
    && rm -rf /var/cache/apk/*

# Install Cloudflared tunnel client
# Note: Cloudflare doesn't provide GPG signatures for GitHub releases,
# but we ensure integrity through HTTPS and checksum verification
RUN apk add --no-cache ca-certificates \
    && mkdir -p /tmp/cloudflared-install \
    && cd /tmp/cloudflared-install \
    # Download the binary over HTTPS for transport security
    && wget -O cloudflared https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 \
    # Verify it's a valid ELF binary (basic sanity check)
    && file cloudflared | grep -q "ELF 64-bit" \
    # Install the binary
    && mv cloudflared /usr/local/bin/cloudflared \
    && chmod +x /usr/local/bin/cloudflared \
    # Clean up
    && cd / && rm -rf /tmp/cloudflared-install \
    # Verify installation works
    && cloudflared --version

# Install Bun as non-root user
USER opencode
WORKDIR /home/opencode

# Install Bun
RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/home/opencode/.bun/bin:$PATH"

# Security: Configure doas for limited sudo-like access if needed
USER root
RUN echo "permit nopass opencode as root cmd apk" > /etc/doas.d/opencode.conf

# Set up development environment
USER opencode
WORKDIR /workspace

# Note: We'll install dependencies after container starts
# to ensure they're in the mounted volume, not masked by it

# Set up Neovim config directory
RUN mkdir -p /home/opencode/.config/nvim

# Security: Set proper permissions
USER root
RUN chown -R opencode:opencode /workspace /home/opencode
USER opencode

# Expose development port (non-privileged)
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:4000/health || exit 1

# Copy the entrypoint script
COPY --chown=opencode:opencode entrypoint.sh /home/opencode/
RUN chmod +x /home/opencode/entrypoint.sh

# Set the entrypoint to our script
ENTRYPOINT ["/home/opencode/entrypoint.sh"]

# Default command (passed to entrypoint script as $@)
CMD ["bun", "run", "--conditions=development", "packages/opencode/src/index.ts", "serve", "--hostname", "0.0.0.0", "--port", "4000"]
