#!/bin/bash

# Setup local-only Go build system without external dependencies

set -e

echo "🔧 Setting up local-only Go build system..."
echo

# 1. Update TUI imports to match actual SDK module
echo "📝 Updating TUI import statements..."
find packages/tui -name "*.go" -exec sed -i 's|github.com/sst/opencode-sdk-go|git.j9xym.com/openapi-api-go|g' {} +

# 2. Update TUI go.mod
echo "📦 Updating TUI go.mod..."
cat > packages/tui/go.mod << 'EOF'
module github.com/sst/opencode

go 1.24.0

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/alecthomas/chroma/v2 v2.18.0
	github.com/charmbracelet/bubbles/v2 v2.0.0-beta.1
	github.com/charmbracelet/bubbletea/v2 v2.0.0-beta.4
	github.com/charmbracelet/glamour v0.10.0
	github.com/charmbracelet/lipgloss/v2 v2.0.0-beta.3
	github.com/charmbracelet/x/ansi v0.9.3
	github.com/fsnotify/fsnotify v1.8.0
	github.com/google/uuid v1.6.0
	github.com/lithammer/fuzzysearch v1.1.8
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6
	github.com/muesli/reflow v0.3.0
	github.com/muesli/termenv v0.16.0
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3
	git.j9xym.com/openapi-api-go v0.0.0
	golang.org/x/image v0.28.0
	rsc.io/qr v0.2.0
)

replace (
	github.com/charmbracelet/x/input => ./input
	git.j9xym.com/openapi-api-go => ../sdk/go
)

require golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
EOF

# 3. Create go.work workspace file
echo "🏗️  Creating Go workspace..."
cat > go.work << 'EOF'
go 1.24.0

toolchain go1.24.4

use (
	./sdk/go
	./packages/tui
	./packages/tui/input
)
EOF

# 4. Initialize SDK module properly
echo "📦 Initializing SDK module..."
cd sdk/go
if [ ! -f "go.mod" ]; then
	cat > go.mod << 'EOF'
module git.j9xym.com/openapi-api-go

go 1.24.4
EOF
fi

# Create basic SDK structure if missing
if [ ! -f "types.gen.go" ]; then
	echo "⚠️  Warning: types.gen.go not found. Run generation script first:"
	echo "   ./scripts/generate-go-sdk.sh --clean --verbose"
fi

cd ../..

echo "✅ Local Go build system setup complete!"
echo
echo "Next steps:"
echo "1. Generate the SDK: ./scripts/generate-go-sdk.sh --clean --verbose"
echo "2. Test the build: go build ./..."
echo "3. Test the TUI: cd packages/tui && go run ."
echo
echo "Benefits:"
echo "• No external GitHub dependencies"
echo "• Local workspace management"
echo "• Consistent module names"
echo "• Proper build isolation"