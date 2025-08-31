#!/bin/bash
set -e  # Exit on any error

echo "🚀 Starting OpenCode development environment..."

# Change to workspace directory
cd /workspace

# Check if we have the necessary package files
if [ ! -f "package.json" ]; then
    echo "❌ Error: package.json not found in /workspace"
    echo "Make sure you're mounting your project directory to /workspace"
    exit 1
fi

echo "🌐 Setting up Cloudflare tunnel for reliable connectivity..."

# Check for Cloudflare tunnel token from various sources
TUNNEL_TOKEN=""

if [ -f "/run/secrets/cloudflare_tunnel_token" ]; then
    echo "✅ Found tunnel token in Podman secret"
    TUNNEL_TOKEN=$(cat /run/secrets/cloudflare_tunnel_token)
elif [ -n "$CLOUDFLARE_TUNNEL_TOKEN" ]; then
    echo "✅ Found tunnel token in environment variable"
    TUNNEL_TOKEN="$CLOUDFLARE_TUNNEL_TOKEN"
elif [ -f "/home/opencode/.credentials/tunnel-token" ]; then
    echo "✅ Found tunnel token in credentials file"
    TUNNEL_TOKEN=$(cat /home/opencode/.credentials/tunnel-token)
else
    echo "⚠️  No Cloudflare tunnel token found, using simplified proxy mode"
    echo "   For full tunnel features, provide token via Podman secret, environment, or credentials file"
fi

if [ -n "$TUNNEL_TOKEN" ]; then
    # Use authenticated tunnel with your token
    echo "🔐 Starting authenticated Cloudflare tunnel..."
    cloudflared tunnel --token "$TUNNEL_TOKEN" --url http://localhost:8080 --no-autoupdate &
else
    # Fall back to simple proxy mode (less reliable but works for basic connectivity)
    echo "🌐 Starting simple tunnel proxy..."
    cloudflared tunnel --config /dev/null --url http://localhost:8080 --no-autoupdate &
fi

# Wait for tunnel to establish
sleep 5

# Configure package managers to use the tunnel proxy
export http_proxy=http://localhost:8080
export https_proxy=http://localhost:8080
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080

echo "✅ Tunnel proxy established on port 8080"

echo "📦 Checking dependencies..."

# Check if node_modules exists and has content
if [ ! -d "node_modules" ] || [ -z "$(ls -A node_modules 2>/dev/null)" ]; then
    echo "📥 Installing dependencies with Bun (via Cloudflare tunnel)..."
    bun install
    echo "✅ Dependencies installed successfully"
else
    echo "📋 Dependencies already present, checking if they're up to date..."
    # Check if package.json or bun.lockb is newer than node_modules
    if [ "package.json" -nt "node_modules" ] || [ "bun.lockb" -nt "node_modules" ]; then
        echo "🔄 Package files are newer than node_modules, updating dependencies..."
        bun install
        echo "✅ Dependencies updated successfully"
    else
        echo "✅ Dependencies are up to date"
    fi
fi

# Verify critical dependencies
echo "🔍 Verifying critical dependencies..."
if bun pm ls | grep -q "zod-openapi"; then
    echo "✅ zod-openapi found"
else
    echo "⚠️ zod-openapi not found, attempting to install..."
    bun add zod-openapi
fi

# Keep tunnel running but don't interfere with main application
echo "🎯 Starting OpenCode server on port 4000 (TUI client port)..."

# Execute the original command passed to the container
exec "$@"
