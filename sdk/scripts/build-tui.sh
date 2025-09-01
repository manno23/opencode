#!/bin/bash

# Build script for the OpenCode TUI binary
# This script builds the TUI binary from packages/tui/cmd/opencode/main.go

set -e

echo "🏗️  Building OpenCode TUI binary..."

# Ensure we're in the right directory
if [ ! -f "go.work" ]; then
    echo "❌ Error: go.work not found. Run from repository root."
    exit 1
fi

# Ensure we're in the repository root
cd "$(dirname "$0")/.."

# Check if the main.go file exists
if [ ! -f "packages/tui/cmd/opencode/main.go" ]; then
    echo "❌ Error: TUI main.go not found at packages/tui/cmd/opencode/main.go"
    exit 1
fi

# Sync workspace
echo "🔄 Syncing workspace..."
go work sync

# Build the TUI binary
echo "🔨 Building TUI binary..."
cd packages/tui

# Tidy dependencies
go mod tidy

# Build the binary
echo "  • Compiling TUI..."
go build -o opencode-tui ./cmd/opencode

# Check if build was successful
if [ ! -f "opencode-tui" ]; then
    echo "❌ Error: Failed to build TUI binary"
    exit 1
fi

echo "✅ TUI binary built successfully: packages/tui/opencode-tui"

# Optional: Show binary info
echo "📊 Binary information:"
file opencode-tui
ls -lh opencode-tui

echo
echo "🎉 TUI build complete!"
echo "Run with: ./packages/tui/opencode-tui"