#!/bin/bash
# =============================================================================
# Go SDK Cleanup Script
# =============================================================================
#
# This script provides comprehensive cleanup of Go SDK artifacts and workspace
# state. It includes automatic backups and selective cleanup options.
#
# Usage:
#   ./scripts/clean-go-sdk.sh [options]
#
# Options:
#   --yes          Skip confirmation prompts
#   --workspace    Clean workspace operations  
#   --vendor       Remove vendor directories
#   --no-backup    Skip automatic backups
#   --deep         Deep clean including module cache
#   --help         Show this help
#
# Example:
#   ./scripts/clean-go-sdk.sh --yes --deep
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
readonly ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
readonly SDK_DIR="$ROOT_DIR/sdk/go"
readonly BACKUP_DIR="$ROOT_DIR/sdk/.go-backup/$(date +%Y%m%d_%H%M%S)"

# Flags
SKIP_CONFIRMATION=false
CLEAN_WORKSPACE=false
REMOVE_VENDOR=false
SKIP_BACKUP=false
DEEP_CLEAN=false

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

confirm() {
    local message="$1"
    if [[ "$SKIP_CONFIRMATION" == "true" ]]; then
        return 0
    fi
    
    read -p "$message (y/N): " -n 1 -r
    echo
    [[ $REPLY =~ ^[Yy]$ ]]
}

run_command() {
    local cmd="$1"
    local description="$2"
    
    log_info "$description"
    if ! eval "$cmd" 2>/dev/null; then
        log_warning "Command failed (non-critical): $description"
    else
        log_success "$description completed"
    fi
}

# =============================================================================
# Backup Functions
# =============================================================================

create_backup() {
    if [[ "$SKIP_BACKUP" == "true" ]]; then
        log_info "Skipping backup (--no-backup specified)"
        return 0
    fi
    
    log_info "Creating backup..."
    mkdir -p "$BACKUP_DIR"
    
    # Files and directories to backup
    local items_to_backup=(
        "sdk/go/types.gen.go"
        "sdk/go/go.mod"
        "sdk/go/go.sum"
        "go.work"
        "go.work.sum"
        "schema/openapi.json"
        "schema/openapi.bundled.json"
    )
    
    for item in "${items_to_backup[@]}"; do
        if [[ -f "$ROOT_DIR/$item" ]]; then
            local backup_path="$BACKUP_DIR/$(dirname "$item")"
            mkdir -p "$backup_path"
            cp "$ROOT_DIR/$item" "$backup_path/"
            log_info "Backed up: $item"
        fi
    done
    
    # Backup generated directories if they exist
    if [[ -d "$SDK_DIR/generated" ]]; then
        cp -r "$SDK_DIR/generated" "$BACKUP_DIR/"
        log_info "Backed up: sdk/go/generated"
    fi
    
    log_success "Backup created at: $BACKUP_DIR"
}

# =============================================================================
# Cleanup Functions  
# =============================================================================

clean_generated_files() {
    log_info "Cleaning generated files..."
    
    # Remove generated Go types
    if [[ -f "$SDK_DIR/types.gen.go" ]]; then
        rm -f "$SDK_DIR/types.gen.go"
        log_success "Removed types.gen.go"
    fi
    
    # Remove generated directories
    if [[ -d "$SDK_DIR/generated" ]]; then
        rm -rf "$SDK_DIR/generated"
        log_success "Removed generated directory"
    fi
    
    # Remove OpenAPI artifacts
    if [[ -f "$ROOT_DIR/schema/openapi.json" ]]; then
        rm -f "$ROOT_DIR/schema/openapi.json"
        log_success "Removed OpenAPI spec"
    fi
    
    if [[ -f "$ROOT_DIR/schema/openapi.bundled.json" ]]; then
        rm -f "$ROOT_DIR/schema/openapi.bundled.json" 
        log_success "Removed bundled OpenAPI spec"
    fi
}

clean_build_artifacts() {
    log_info "Cleaning build artifacts..."
    
    cd "$ROOT_DIR"
    
    # Remove common build directories
    local build_dirs=("bin" "dist" "build" "tmp" ".tmp")
    for dir in "${build_dirs[@]}"; do
        if [[ -d "$dir" ]]; then
            rm -rf "$dir"
            log_success "Removed $dir directory"
        fi
    done
    
    # Remove binary files
    find . -maxdepth 3 -name '*.exe' -o -name '*.dll' -o -name '*.so' -o -name '*.dylib' | while read -r file; do
        rm -f "$file"
        log_success "Removed binary: $file"
    done
}

clean_go_modules() {
    log_info "Cleaning Go module files..."
    
    cd "$SDK_DIR"
    
    # Remove go.sum to force re-verification
    if [[ -f "go.sum" ]]; then
        rm -f go.sum
        log_success "Removed go.sum"
    fi
    
    # Clean workspace sum if requested
    if [[ "$CLEAN_WORKSPACE" == "true" && -f "$ROOT_DIR/go.work.sum" ]]; then
        rm -f "$ROOT_DIR/go.work.sum"
        log_success "Removed go.work.sum"
    fi
}

clean_vendor_directories() {
    if [[ "$REMOVE_VENDOR" == "false" ]]; then
        return 0
    fi
    
    log_info "Cleaning vendor directories..."
    
    # Remove SDK vendor directory
    if [[ -d "$SDK_DIR/vendor" ]]; then
        rm -rf "$SDK_DIR/vendor"
        log_success "Removed SDK vendor directory"
    fi
    
    # Remove TUI vendor directory
    local tui_vendor="$ROOT_DIR/packages/tui/vendor"
    if [[ -d "$tui_vendor" ]]; then
        rm -rf "$tui_vendor"
        log_success "Removed TUI vendor directory"
    fi
}

clean_go_caches() {
    if [[ "$DEEP_CLEAN" == "false" ]]; then
        return 0
    fi
    
    log_info "Performing deep clean of Go caches..."
    
    cd "$ROOT_DIR"
    
    # Clean all Go caches
    run_command "go clean -cache" "Cleaning build cache"
    run_command "go clean -testcache" "Cleaning test cache"
    run_command "go clean -modcache" "Cleaning module cache"
    
    # Clean additional cache directories
    if [[ -d "$HOME/.cache/go-build" ]]; then
        rm -rf "$HOME/.cache/go-build"
        log_success "Removed user Go build cache"
    fi
}

reset_workspace() {
    if [[ "$CLEAN_WORKSPACE" == "false" ]]; then
        return 0
    fi
    
    log_info "Resetting workspace state..."
    
    cd "$ROOT_DIR"
    
    if [[ -f "go.work" ]]; then
        run_command "go work sync" "Syncing workspace"
        log_success "Workspace reset completed"
    else
        log_warning "No go.work file found"
    fi
}

# =============================================================================
# Restoration Functions
# =============================================================================

restore_from_backup() {
    local backup_pattern=".go-cleanup-backup-*"
    
    # Find the most recent backup
    local latest_backup
    latest_backup=$(find "$ROOT_DIR" -maxdepth 1 -name "$backup_pattern" -type d | sort | tail -1)
    
    if [[ -z "$latest_backup" ]]; then
        log_warning "No backup found for restoration"
        return 1
    fi
    
    log_info "Found backup: $latest_backup"
    
    if confirm "Restore from backup?"; then
        log_info "Restoring from backup..."
        
        # Restore key files
        if [[ -f "$latest_backup/go.mod" ]]; then
            cp "$latest_backup/go.mod" "$SDK_DIR/"
            log_success "Restored go.mod"
        fi
        
        if [[ -f "$latest_backup/go.sum" ]]; then
            cp "$latest_backup/go.sum" "$SDK_DIR/"
            log_success "Restored go.sum"
        fi
        
        if [[ -f "$latest_backup/types.gen.go" ]]; then
            cp "$latest_backup/types.gen.go" "$SDK_DIR/"
            log_success "Restored types.gen.go"
        fi
        
        log_success "Backup restoration completed"
        return 0
    fi
    
    return 1
}

# =============================================================================
# Help and Argument Parsing
# =============================================================================

show_help() {
    cat << EOF
Go SDK Cleanup Script

This script provides comprehensive cleanup of Go SDK artifacts and workspace
state with automatic backup capabilities.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --yes           Skip confirmation prompts
    --workspace     Clean workspace operations (go.work.sum, sync)
    --vendor        Remove vendor directories
    --no-backup     Skip automatic backups (use with caution)
    --deep          Deep clean including module cache
    --help          Show this help message

EXAMPLES:
    $0                       # Standard cleanup with prompts
    $0 --yes --deep          # Deep clean without prompts
    $0 --workspace --vendor  # Clean workspace and vendor dirs

CLEANUP OPERATIONS:
    - Remove generated types (types.gen.go)
    - Remove generated directories
    - Remove OpenAPI specifications  
    - Remove build artifacts
    - Clean Go module files (go.sum)
    - [--vendor] Remove vendor directories
    - [--workspace] Clean workspace files
    - [--deep] Clean Go caches and module cache

BACKUP:
    By default, important files are backed up before cleanup:
    - types.gen.go
    - go.mod, go.sum  
    - go.work, go.work.sum
    - OpenAPI specifications
    
RESTORATION:
    If cleanup causes issues, restore from backup:
    cp .go-cleanup-backup-*/types.gen.go sdk/go/
    cp .go-cleanup-backup-*/go.sum sdk/go/

EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --yes)
                SKIP_CONFIRMATION=true
                shift
                ;;
            --workspace)
                CLEAN_WORKSPACE=true
                shift
                ;;
            --vendor)
                REMOVE_VENDOR=true
                shift
                ;;
            --no-backup)
                SKIP_BACKUP=true
                shift
                ;;
            --deep)
                DEEP_CLEAN=true
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
# Main Function
# =============================================================================

main() {
    echo
    echo "========================================"
    echo "  Go SDK Cleanup Script"
    echo "========================================"
    echo
    
    # Parse arguments
    parse_arguments "$@"
    
    # Change to root directory
    cd "$ROOT_DIR"
    
    # Safety confirmation
    if ! confirm "This will clean Go SDK artifacts and workspace state. Continue?"; then
        log_info "Cleanup cancelled by user"
        exit 0
    fi
    
    # Show what will be cleaned
    echo
    log_info "Cleanup operations selected:"
    echo "  - Generated files: YES"
    echo "  - Build artifacts: YES"
    echo "  - Go modules: YES"
    echo "  - Vendor directories: $([ "$REMOVE_VENDOR" = true ] && echo "YES" || echo "NO")"
    echo "  - Workspace: $([ "$CLEAN_WORKSPACE" = true ] && echo "YES" || echo "NO")"
    echo "  - Deep clean (caches): $([ "$DEEP_CLEAN" = true ] && echo "YES" || echo "NO")"
    echo "  - Create backup: $([ "$SKIP_BACKUP" = true ] && echo "NO" || echo "YES")"
    echo
    
    # Create backup first
    create_backup
    
    # Perform cleanup operations
    clean_generated_files
    clean_build_artifacts  
    clean_go_modules
    clean_vendor_directories
    clean_go_caches
    reset_workspace
    
    # Success message
    echo
    log_success "Go SDK cleanup completed successfully!"
    
    if [[ "$SKIP_BACKUP" == "false" ]]; then
        echo
        log_info "Backup available at: $BACKUP_DIR"
        log_info "To restore if needed:"
        echo "  cp $BACKUP_DIR/types.gen.go sdk/go/"
        echo "  cp $BACKUP_DIR/go.sum sdk/go/"
    fi
    
    echo
    log_info "Next steps:"
    echo "  1. Run './scripts/generate-go-sdk.sh' to regenerate SDK"
    echo "  2. Test build with 'go build ./...' from project root" 
    echo "  3. Verify TUI integration still works"
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

# Handle special case: restoration mode
if [[ $# -eq 1 && "$1" == "--restore" ]]; then
    restore_from_backup
    exit $?
fi

# Run main function
main "$@"
