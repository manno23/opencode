# Use minimal Alpine base for security
FROM node:20-alpine

# Security: Create non-root user
RUN addgroup -g 1001 -S opencode && \
    adduser -S opencode -u 1001 -G opencode

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
