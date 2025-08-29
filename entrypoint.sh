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

echo "📦 Checking dependencies..."

# Check if node_modules exists and has content
if [ ! -d "node_modules" ] || [ -z "$(ls -A node_modules 2>/dev/null)" ]; then
    echo "📥 Installing dependencies with Bun..."
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

echo "🎯 Starting OpenCode server..."

# Execute the original command passed to the container
exec "$@"
