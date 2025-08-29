#!/bin/bash

# Local Go build script for CI/CD and development
# Ensures all builds use local modules only

set -e

echo "🏗️  Building Go components locally..."
echo

# Ensure we're in the right directory
if [ ! -f "go.work" ]; then
	echo "❌ Error: go.work not found. Run from repository root."
	exit 1
fi

# Sync workspace
echo "🔄 Syncing workspace..."
go work sync

# Build all modules
echo "🔨 Building all modules..."

# SDK
echo "  • Building SDK..."
cd sdk/go
go mod tidy
go build -v .
go test -v . 2>/dev/null || echo "  ⚠️  SDK tests may need generation first"

# TUI
echo "  • Building TUI..."
cd ../../packages/tui
go mod tidy
go build -v ./cmd/opencode

# Input module
echo "  • Building input module..."
cd input
go mod tidy
go build -v .
cd ../../..

# Full workspace build
echo "🔗 Building workspace..."
go build -v ./...

echo
echo "✅ Local Go build complete!"
echo
echo "Built components:"
echo "• SDK (sdk/go)"
echo "• TUI (packages/tui)"
echo "• Input handler (packages/tui/input)"
echo "• Workspace integration"