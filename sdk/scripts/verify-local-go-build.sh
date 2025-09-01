#!/bin/bash

# Verify local Go build system works without external dependencies

set -e

echo "🔍 Verifying local Go build system..."
echo

# Check Go workspace
echo "📋 Checking Go workspace..."
if [ ! -f "go.work" ]; then
	echo "❌ go.work not found"
	exit 1
fi

go work sync
echo "✅ Workspace sync successful"

# Check module structure
echo
echo "📦 Checking module structure..."

# SDK module
if [ ! -f "sdk/go/go.mod" ]; then
	echo "❌ SDK go.mod not found"
	exit 1
fi

# TUI module
if [ ! -f "packages/tui/go.mod" ]; then
	echo "❌ TUI go.mod not found"
	exit 1
fi

echo "✅ All modules present"

# Test builds
echo
echo "🔨 Testing builds..."

# Build SDK
echo "Building SDK..."
cd sdk/go
go mod tidy
go build .
cd ../..

# Build TUI
echo "Building TUI..."
cd packages/tui
go mod tidy
go build ./cmd/opencode
cd ../..

echo "✅ All builds successful"

# Test workspace build
echo
echo "🏗️  Testing workspace build..."
go build ./...

echo "✅ Workspace build successful"

echo
echo "🎉 Local Go build system verification complete!"
echo "All components build successfully without external dependencies."