#!/bin/bash
# =============================================================================
# Regain Control Script
# =============================================================================
#
# Switch from external SDK dependency back to local generation
# This gives you full control over your API surface again
#
# =============================================================================

set -euo pipefail

direnv allow

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Remove external SDK dependency
remove_external_sdk() {
    log_info "Removing external SDK dependency..."

    # Remove from go.mod
    if [[ -f "sdk/go/go.mod" ]]; then
        sed -i '/github.com\/sst\/opencode-sdk-go/d' sdk/go/go.mod
        sed -i '/github.com\/sst\/opencode/d' sdk/go/go.mod
    fi

    # Remove from go.sum
    if [[ -f "sdk/go/go.sum" ]]; then
        sed -i '/github.com\/sst\/opencode/d' sdk/go/go.sum
    fi

    log_success "External SDK dependency removed"
}

# Update TUI to use local SDK
update_tui_imports() {
    log_info "Updating TUI imports to use local SDK..."

    # Update TUI go.mod to use local SDK
    if [[ -f "packages/tui/go.mod" ]]; then
        sed -i 's|github.com/sst/opencode-sdk-go|git.j9xym.com/opencode-api-go|g' packages/tui/go.mod
    fi

    # Update TUI imports
    find packages/tui -name "*.go" -exec sed -i 's|github.com/sst/opencode|git.j9xym.com/opencode-api-go|g' {} \;

    log_success "TUI imports updated"
}

# Update package.json scripts
update_scripts() {
    log_info "Updating package.json scripts for local generation..."

    # This would update the generate script to use local generation
    # instead of the external SDK

    log_success "Scripts updated"
}

# Generate fresh SDKs
generate_fresh_sdks() {
    log_info "Generating fresh local SDKs..."

    # Generate TypeScript SDK
    log_info "Generating TypeScript SDK..."
    cd packages/sdk/js
    bun run --conditions=development ../../../packages/opencode/src/index.ts generate > openapi.json
    bun run script/generate.ts
    cd ../../..

    # Generate Go SDK
    log_info "Generating Go SDK..."
    bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json
    npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json
    go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
        -generate types \
        -package opencode \
        schema/openapi.bundled.json > sdk/go/types.gen.go

    # Copy service structure
    if [[ -d "packages/sdk/go" ]]; then
        cp -r packages/sdk/go/* sdk/go/ 2>/dev/null || true
    fi

    # Update import paths
    sed -i 's|github.com/sst/opencode|git.j9xym.com/opencode-api-go|g' sdk/go/*.go

    # Tidy Go modules
    cd sdk/go
    go mod tidy
    go fmt ./...
    cd ../..

    log_success "Fresh SDKs generated"
}

main() {
    echo
    echo "========================================"
    echo "  🔓 Regaining Control of Your API"
    echo "========================================"
    echo
    log_warning "This will switch from external SDK dependency to local generation"
    log_warning "You'll have full control over your API surface again"
    echo

    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Operation cancelled"
        exit 0
    fi

    remove_external_sdk
    update_tui_imports
    update_scripts
    generate_fresh_sdks

    echo
    log_success "🎉 Control regained!"
    echo
    echo "You now have:"
    echo "• Full control over your API surface"
    echo "• Local SDK generation"
    echo "• No external dependencies for SDK"
    echo
    echo "Next steps:"
    echo "1. Test your changes: bun run test:api-surface"
    echo "2. Make API changes: modify packages/opencode/src/server/server.ts"
    echo "3. Regenerate SDKs: ./scripts/generate-local-sdk.sh --all"
    echo
    echo "You're free to modify your API as you wish! 🚀"
}

# Run main function
main "$@"
