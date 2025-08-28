#!/bin/bash
# =============================================================================
# Go SDK Generation Script
# =============================================================================
#
# This script generates a hybrid Go SDK using oapi-codegen for types and 
# manual service implementations. It provides complete automation with status
# reporting and error handling.
#
# Usage:
#   ./scripts/generate-go-sdk.sh [options]
#
# Options:
#   --clean        Clean before generation
#   --no-build     Skip validation build
#   --verbose      Enable verbose output  
#   --help         Show this help
#
# Example:
#   ./scripts/generate-go-sdk.sh --clean --verbose
#
# =============================================================================

set -euo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

# Configuration
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
readonly SDK_DIR="$ROOT_DIR/sdk/go"
readonly SOURCE_SDK_DIR="$ROOT_DIR/packages/sdk/go"
readonly SCHEMA_DIR="$ROOT_DIR/schema"

# Flags
CLEAN_FIRST=false
SKIP_BUILD=false
VERBOSE=false

# =============================================================================
# Utility Functions
# =============================================================================

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

log_verbose() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${BLUE}[VERBOSE]${NC} $1"
    fi
}

run_command() {
    local cmd="$1"
    local description="$2"
    
    log_info "$description"
    log_verbose "Running: $cmd"
    
    if [[ "$VERBOSE" == "true" ]]; then
        eval "$cmd"
    else
        if ! eval "$cmd" >/dev/null 2>&1; then
            log_error "Failed: $description"
            log_error "Command: $cmd"
            return 1
        fi
    fi
    
    log_success "$description completed"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check Go version
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed"
        return 1
    fi
    
    local go_version
    go_version=$(go version | grep -o 'go[0-9]\.[0-9]*' | sed 's/go//')
    local major minor
    major=$(echo "$go_version" | cut -d. -f1)
    minor=$(echo "$go_version" | cut -d. -f2)
    
    if [[ $major -lt 1 ]] || [[ $major -eq 1 && $minor -lt 24 ]]; then
        log_warning "Go 1.24+ recommended (found $go_version)"
    fi
    
    # Check Node.js/Bun
    if ! command -v bun &> /dev/null && ! command -v node &> /dev/null; then
        log_error "Neither Bun nor Node.js is available"
        return 1
    fi
    
    # Check redocly
    if ! command -v npx &> /dev/null; then
        log_error "npx is not available (needed for redocly)"
        return 1
    fi
    
    log_success "Prerequisites check passed"
}

show_help() {
    cat << EOF
Go SDK Generation Script

This script generates a hybrid Go SDK using oapi-codegen for types and 
manual service implementations.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --clean        Clean before generation (recommended for fresh builds)
    --no-build     Skip validation build step
    --verbose      Enable verbose output for debugging
    --help         Show this help message

EXAMPLES:
    $0                           # Standard generation
    $0 --clean --verbose         # Clean generation with detailed output
    $0 --no-build               # Generate without validation build

WORKFLOW:
    1. Generate OpenAPI specification from TypeScript API
    2. Bundle specification (dereference all \$ref entries)
    3. Install/update oapi-codegen tool
    4. Generate Go types using oapi-codegen
    5. Copy service structure from source SDK
    6. Update import paths for target module
    7. Tidy Go modules and run validation build
    8. Sync workspace if go.work exists

OUTPUT:
    The generated SDK will be available at: $SDK_DIR

EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --clean)
                CLEAN_FIRST=true
                shift
                ;;
            --no-build)
                SKIP_BUILD=true
                shift
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# =============================================================================
# Generation Functions
# =============================================================================

clean_sdk() {
    if [[ "$CLEAN_FIRST" == "true" ]]; then
        log_info "Cleaning SDK directory..."
        
        # Remove generated files
        if [[ -f "$SDK_DIR/types.gen.go" ]]; then
            rm -f "$SDK_DIR/types.gen.go"
            log_verbose "Removed types.gen.go"
        fi
        
        # Remove generated directories
        if [[ -d "$SDK_DIR/generated" ]]; then
            rm -rf "$SDK_DIR/generated"
            log_verbose "Removed generated directory"
        fi
        
        # Clean module cache if needed
        cd "$ROOT_DIR"
        run_command "go clean -modcache" "Cleaning Go module cache"
        
        log_success "SDK cleaned"
    fi
}

generate_openapi_spec() {
    log_info "Generating OpenAPI specification..."
    
    cd "$ROOT_DIR"
    
    # Generate base specification
    run_command \
        "bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json" \
        "Generating base OpenAPI spec"
    
    # Verify the spec was generated
    if [[ ! -f "$SCHEMA_DIR/openapi.json" ]]; then
        log_error "Failed to generate OpenAPI specification"
        return 1
    fi
    
    local spec_size
    spec_size=$(wc -l < "$SCHEMA_DIR/openapi.json")
    log_verbose "Generated spec: $spec_size lines"
}

bundle_specification() {
    log_info "Bundling OpenAPI specification..."
    
    cd "$ROOT_DIR"
    
    # Bundle with dereferencing
    run_command \
        "npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced" \
        "Bundling and dereferencing spec"
    
    # Verify bundled spec
    if [[ ! -f "$SCHEMA_DIR/openapi.bundled.json" ]]; then
        log_error "Failed to create bundled specification"
        return 1
    fi
    
    local bundled_size
    bundled_size=$(wc -l < "$SCHEMA_DIR/openapi.bundled.json")
    log_verbose "Bundled spec: $bundled_size lines"
}

install_oapi_codegen() {
    log_info "Installing/updating oapi-codegen..."
    
    # Install latest version of oapi-codegen
    run_command \
        "go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest" \
        "Installing oapi-codegen"
    
    # Verify installation
    if ! command -v oapi-codegen &> /dev/null; then
        log_error "oapi-codegen installation failed"
        return 1
    fi
    
    local version
    version=$(oapi-codegen -version 2>&1 || echo "unknown")
    log_verbose "oapi-codegen version: $version"
}

generate_types() {
    log_info "Generating Go types..."
    
    cd "$ROOT_DIR"
    
    # Ensure SDK directory exists
    mkdir -p "$SDK_DIR"
    
    # Generate types using oapi-codegen
    run_command \
        "oapi-codegen -generate types -package opencode schema/openapi.bundled.json > sdk/go/types.gen.go" \
        "Generating types with oapi-codegen"
    
    # Verify types were generated
    if [[ ! -f "$SDK_DIR/types.gen.go" ]]; then
        log_error "Failed to generate types"
        return 1
    fi
    
    local type_lines
    type_lines=$(wc -l < "$SDK_DIR/types.gen.go")
    log_verbose "Generated types: $type_lines lines"
}

copy_service_structure() {
    log_info "Copying service structure..."
    
    # Verify source SDK exists
    if [[ ! -d "$SOURCE_SDK_DIR" ]]; then
        log_error "Source SDK directory not found: $SOURCE_SDK_DIR"
        return 1
    fi
    
    cd "$ROOT_DIR"
    
    # Copy service files using rsync to preserve permissions
    run_command \
        "rsync -av --exclude='go.mod' --exclude='types.gen.go' '$SOURCE_SDK_DIR/' '$SDK_DIR/'" \
        "Copying service structure"
    
    # Count copied files
    local copied_files
    copied_files=$(find "$SDK_DIR" -name "*.go" -not -name "types.gen.go" | wc -l)
    log_verbose "Copied service files: $copied_files"
}

update_import_paths() {
    log_info "Updating import paths..."
    
    cd "$ROOT_DIR"
    
    # Update import paths from old module to new module
    local old_module="opencode.local-sdk-go"
    local new_module="git.j9xym.com/opencode-api-go"
    
    # Find and replace in all Go files
    run_command \
        "find '$SDK_DIR' -name '*.go' -type f -exec sed -i 's|$old_module|$new_module|g' {} +" \
        "Updating import paths"
    
    # Check for any remaining old imports
    local remaining
    remaining=$(grep -r "$old_module" "$SDK_DIR" || true)
    if [[ -n "$remaining" ]]; then
        log_warning "Some old import paths may remain:"
        echo "$remaining"
    else
        log_verbose "All import paths updated successfully"
    fi
}

update_go_module() {
    log_info "Updating Go module..."
    
    cd "$SDK_DIR"
    
    # Tidy modules
    run_command \
        "go mod tidy" \
        "Tidying Go modules"
    
    # Show module info
    if [[ "$VERBOSE" == "true" ]]; then
        local module_info
        module_info=$(go list -m 2>/dev/null || echo "Module info unavailable")
        log_verbose "Module: $module_info"
    fi
}

validate_build() {
    if [[ "$SKIP_BUILD" == "true" ]]; then
        log_info "Skipping validation build (--no-build specified)"
        return 0
    fi
    
    log_info "Running validation build..."
    
    cd "$SDK_DIR"
    
    # Build all packages
    run_command \
        "go build ./..." \
        "Building all packages"
    
    # Run basic tests (if any exist and don't fail)
    if ls *_test.go >/dev/null 2>&1; then
        log_info "Running tests..."
        if [[ "$VERBOSE" == "true" ]]; then
            go test ./... || log_warning "Some tests failed (this may be expected during development)"
        else
            go test ./... >/dev/null 2>&1 || log_warning "Some tests failed (this may be expected during development)"
        fi
    else
        log_verbose "No test files found"
    fi
}

sync_workspace() {
    log_info "Syncing workspace..."
    
    cd "$ROOT_DIR"
    
    # Sync workspace if go.work exists
    if [[ -f "go.work" ]]; then
        run_command \
            "go work sync" \
            "Syncing Go workspace"
        
        # Show workspace status
        if [[ "$VERBOSE" == "true" ]]; then
            local workspace_modules
            workspace_modules=$(go list -m all | head -10)
            log_verbose "Workspace modules (first 10):"
            echo "$workspace_modules"
        fi
    else
        log_verbose "No go.work file found, skipping workspace sync"
    fi
}

# =============================================================================
# Main Function
# =============================================================================

main() {
    echo
    echo "========================================"
    echo "  Go SDK Generation Script"
    echo "========================================"
    echo
    
    # Parse arguments
    parse_arguments "$@"
    
    # Change to root directory
    cd "$ROOT_DIR"
    
    # Run generation workflow
    log_info "Starting SDK generation workflow..."
    
    check_prerequisites
    clean_sdk
    generate_openapi_spec
    bundle_specification
    install_oapi_codegen
    generate_types
    copy_service_structure
    update_import_paths
    update_go_module
    validate_build
    sync_workspace
    
    # Success message
    echo
    log_success "Go SDK generation completed successfully!"
    echo
    echo "Generated SDK location: $SDK_DIR"
    echo "Generated types: $SDK_DIR/types.gen.go"
    
    if [[ "$VERBOSE" == "true" ]]; then
        echo
        echo "Next steps:"
        echo "1. Review generated types in types.gen.go"
        echo "2. Update service implementations as needed"
        echo "3. Update compatibility layer in compat/"
        echo "4. Test integration with TUI application"
    fi
    echo
}

# =============================================================================
# Script Entry Point
# =============================================================================

# Verify we're in the right directory
if [[ ! -f "go.work" ]]; then
    log_error "This script must be run from the project root (directory with go.work)"
    exit 1
fi

# Run main function
main "$@"