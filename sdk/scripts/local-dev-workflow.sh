#!/bin/bash
# =============================================================================
# Local Development Workflow for SDK Generation
# =============================================================================
#
# Streamlined workflow for local development and testing of SDK generation
# Provides fast feedback and comprehensive validation
#
# Usage:
#   ./scripts/local-dev-workflow.sh [command] [options]
#
# Commands:
#   generate      Generate SDK with validation
#   test          Run test suites
#   validate      Run quality validation
#   clean         Clean and regenerate
#   watch         Watch for changes and auto-regenerate
#   ci            Simulate CI pipeline locally
#
# Options:
#   --verbose     Enable verbose output
#   --no-build    Skip build validation
#   --fast        Skip slow validations
#   --help        Show this help
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

# Flags
VERBOSE=false
NO_BUILD=false
FAST=false

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

check_dependencies() {
    local missing_deps=()

    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    fi

    if ! command -v npm && ! command -v bun &> /dev/null&> /dev/null; then
        missing_deps+=("bun or npm")
    fi

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        log_info "Please install missing dependencies and try again"
        exit 1
    fi
}

run_with_timer() {
    local cmd="$1"
    local description="$2"

    log_info "$description..."
    local start_time=$(date +%s)

    if eval "$cmd"; then
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        log_success "$description completed in ${duration}s"
    else
        log_error "$description failed"
        return 1
    fi
}

# =============================================================================
# Workflow Commands
# =============================================================================

cmd_generate() {
    log_info "Starting SDK generation workflow..."

    # Install dependencies
    run_with_timer "bun install" "Installing dependencies"

    # Generate SDK
    local flags=""
    if [[ "$VERBOSE" == "true" ]]; then
        flags="$flags --verbose"
    fi
    if [[ "$NO_BUILD" == "true" ]]; then
        flags="$flags --no-build"
    fi

    run_with_timer "./scripts/generate-go-sdk.sh $flags" "Generating Go SDK"

    # Basic validation
    if [[ -f "sdk/go/types.gen.go" ]]; then
        log_success "SDK generated successfully"
        echo "Location: sdk/go/"
        echo "Types: sdk/go/types.gen.go"
    else
        log_error "SDK generation failed"
        return 1
    fi
}

cmd_test() {
    log_info "Running test suites..."

    # Unit tests
    if [[ -x "./scripts/test-sdk-generation.sh" ]]; then
        local flags=""
        if [[ "$VERBOSE" == "true" ]]; then
            flags="$flags --verbose"
        fi
        if [[ "$FAST" == "true" ]]; then
            flags="$flags --unit --integration"
        else
            flags="$flags --all"
        fi

        run_with_timer "./scripts/test-sdk-generation.sh $flags" "Running SDK generation tests"
    fi

    # Go tests if SDK exists
    if [[ -d "sdk/go" ]]; then
        cd sdk/go
        if ls *_test.go >/dev/null 2>&1; then
            run_with_timer "go test -v ./..." "Running Go unit tests"
        else
            log_info "No Go test files found"
        fi
        cd "$ROOT_DIR"
    fi
}

cmd_validate() {
    log_info "Running quality validation..."

    if [[ ! -x "./scripts/validate-sdk-quality.sh" ]]; then
        log_error "Validation script not found"
        return 1
    fi

    local flags=""
    if [[ "$VERBOSE" == "true" ]]; then
        flags="$flags --verbose"
    fi
    if [[ "$FAST" == "true" ]]; then
        flags="$flags --build --format"
    else
        flags="$flags --all"
    fi

    run_with_timer "./scripts/validate-sdk-quality.sh $flags" "Validating SDK quality"
}

cmd_clean() {
    log_info "Clean regeneration workflow..."

    # Clean existing SDK
    if [[ -d "sdk/go" ]]; then
        run_with_timer "rm -rf sdk/go" "Removing existing SDK"
    fi

    # Clean schema files
    if [[ -f "schema/openapi.bundled.json" ]]; then
        run_with_timer "rm -f schema/openapi.bundled.json" "Cleaning bundled schema"
    fi

    # Regenerate
    cmd_generate
}

cmd_watch() {
    log_info "Starting watch mode..."
    log_info "Watching for changes in packages/opencode/src/, schema/, and scripts/"
    log_info "Press Ctrl+C to stop"

    # Check if fswatch is available
    if ! command -v fswatch &> /dev/null; then
        log_error "fswatch not available. Install with: brew install fswatch"
        log_info "Falling back to manual regeneration..."
        cmd_generate
        exit 0
    fi

    # Watch for changes
    fswatch -o packages/opencode/src/ schema/ scripts/generate-go-sdk.sh | while read -r num; do
        log_info "Changes detected, regenerating SDK..."
        cmd_generate || log_warning "Regeneration failed, but continuing to watch..."
    done
}

cmd_ci() {
    log_info "Simulating CI pipeline locally..."

    # Step 1: Clean and generate
    cmd_clean

    # Step 2: Run tests
    cmd_test

    # Step 3: Validate quality
    cmd_validate

    # Step 4: Build validation
    if [[ "$NO_BUILD" == "false" ]]; then
        run_with_timer "./scripts/generate-go-sdk.sh --no-build" "Final build validation"
    fi

    log_success "CI simulation completed!"
}

# =============================================================================
# Main Function
# =============================================================================

show_help() {
    cat << EOF
Local Development Workflow for SDK Generation

This script provides streamlined workflows for local development and testing.

USAGE:
    $0 COMMAND [OPTIONS]

COMMANDS:
    generate      Generate SDK with basic validation
    test          Run comprehensive test suites
    validate      Run quality validation checks
    clean         Clean existing SDK and regenerate
    watch         Watch for changes and auto-regenerate
    ci            Simulate full CI pipeline locally

OPTIONS:
    --verbose     Enable verbose output for debugging
    --no-build    Skip build validation steps
    --fast        Skip slow validations for faster feedback
    --help        Show this help message

EXAMPLES:
    $0 generate                    # Quick SDK generation
    $0 test --verbose             # Run tests with verbose output
    $0 validate --fast            # Quick validation
    $0 clean                      # Clean and regenerate
    $0 watch                      # Auto-regenerate on changes
    $0 ci --no-build              # Simulate CI without build

WORKFLOW RECOMMENDATIONS:
    Development:  $0 watch
    Testing:      $0 test
    Validation:   $0 validate
    CI Check:     $0 ci

EOF
}

parse_arguments() {
    if [[ $# -eq 0 ]]; then
        show_help
        exit 0
    fi

    COMMAND="$1"
    shift

    while [[ $# -gt 0 ]]; do
        case $1 in
            --verbose)
                VERBOSE=true
                shift
                ;;
            --no-build)
                NO_BUILD=true
                shift
                ;;
            --fast)
                FAST=true
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

main() {
    echo
    echo "========================================"
    echo "  Local SDK Development Workflow"
    echo "========================================"
    echo

    parse_arguments "$@"

    check_dependencies

    cd "$ROOT_DIR"

    case "$COMMAND" in
        generate)
            cmd_generate
            ;;
        test)
            cmd_test
            ;;
        validate)
            cmd_validate
            ;;
        clean)
            cmd_clean
            ;;
        watch)
            cmd_watch
            ;;
        ci)
            cmd_ci
            ;;
        *)
            log_error "Unknown command: $COMMAND"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
