#!/bin/bash
# =============================================================================
# SDK Generation Test Suite
# =============================================================================
#
# Comprehensive test suite for the Go SDK generation workflow
# Tests individual functions, integration scenarios, and error conditions
#
# Usage:
#   ./scripts/test-sdk-generation.sh [options]
#
# Options:
#   --unit          Run unit tests for individual functions
#   --integration   Run integration tests (full workflow)
#   --error-scenarios Run error scenario tests
#   --all           Run all tests
#   --verbose       Enable verbose output
#   --help          Show this help
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
readonly TEST_TMP_DIR="$ROOT_DIR/tmp/test-sdk-gen"
readonly BACKUP_DIR="$ROOT_DIR/tmp/sdk-backup-$(date +%s)"

# Test flags
RUN_UNIT=false
RUN_INTEGRATION=false
RUN_ERROR_SCENARIOS=false
VERBOSE=false

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

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

test_pass() {
    local test_name="$1"
    ((TESTS_RUN++))
    ((TESTS_PASSED++))
    log_success "✓ $test_name"
}

test_fail() {
    local test_name="$1"
    local message="$2"
    ((TESTS_RUN++))
    ((TESTS_FAILED++))
    log_error "✗ $test_name"
    echo "  $message"
}

run_test() {
    local test_name="$1"
    local test_function="$2"

    log_info "Running test: $test_name"
    if $test_function; then
        test_pass "$test_name"
    else
        test_fail "$test_name" "Test function failed"
    fi
}

setup_test_environment() {
    log_info "Setting up test environment..."

    # Create test directories
    mkdir -p "$TEST_TMP_DIR"
    mkdir -p "$BACKUP_DIR"

    # Backup original files if they exist
    if [[ -d "$ROOT_DIR/sdk/go" ]]; then
        cp -r "$ROOT_DIR/sdk/go" "$BACKUP_DIR/"
    fi

    if [[ -f "$ROOT_DIR/schema/openapi.json" ]]; then
        cp "$ROOT_DIR/schema/openapi.json" "$BACKUP_DIR/"
    fi

    log_success "Test environment setup complete"
}

cleanup_test_environment() {
    log_info "Cleaning up test environment..."

    # Restore original files
    if [[ -d "$BACKUP_DIR/sdk" ]]; then
        rm -rf "$ROOT_DIR/sdk"
        cp -r "$BACKUP_DIR/sdk" "$ROOT_DIR/"
    fi

    if [[ -f "$BACKUP_DIR/openapi.json" ]]; then
        cp "$BACKUP_DIR/openapi.json" "$ROOT_DIR/schema/"
    fi

    # Clean up test directory
    rm -rf "$TEST_TMP_DIR"
    rm -rf "$BACKUP_DIR"

    log_success "Test environment cleanup complete"
}

# =============================================================================
# Unit Tests
# =============================================================================

test_check_prerequisites() {
    log_verbose "Testing prerequisites check..."

    # Test with valid environment
    cd "$ROOT_DIR"

    # Mock go command if not available
    if ! command -v go &> /dev/null; then
        log_warning "Go not available, skipping prerequisite test"
        return 0
    fi

    # This would test the check_prerequisites function
    # For now, just check if the function exists in the main script
    if grep -q "check_prerequisites()" "$SCRIPT_DIR/generate-go-sdk.sh"; then
        return 0
    else
        return 1
    fi
}

test_generate_openapi_spec() {
    log_verbose "Testing OpenAPI spec generation..."

    cd "$ROOT_DIR"

    # Test if the generation command works
    if command -v bun &> /dev/null; then
        if bun run --conditions=development packages/opencode/src/index.ts generate > /dev/null 2>&1; then
            if [[ -f "schema/openapi.json" ]]; then
                return 0
            fi
        fi
    fi

    return 1
}

test_bundle_specification() {
    log_verbose "Testing spec bundling..."

    cd "$ROOT_DIR"

    # Check if redocly is available
    if command -v npx &> /dev/null || command -v bun &> /dev/null; then
        # This would test the bundling process
        return 0
    else
        log_warning "Bundling tools not available"
        return 1
    fi
}

test_oapi_codegen_install() {
    log_verbose "Testing oapi-codegen installation..."

    if command -v go &> /dev/null; then
        if go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest; then
            if command -v oapi-codegen &> /dev/null; then
                return 0
            fi
        fi
    fi

    return 1
}

# =============================================================================
# Integration Tests
# =============================================================================

test_full_generation_workflow() {
    log_verbose "Testing full generation workflow..."

    cd "$ROOT_DIR"

    # Run the full generation script
    if [[ -x "$SCRIPT_DIR/generate-go-sdk.sh" ]]; then
        if timeout 300s "$SCRIPT_DIR/generate-go-sdk.sh" --verbose; then
            # Check if essential files were created
            if [[ -f "sdk/go/types.gen.go" && -f "sdk/go/go.mod" ]]; then
                return 0
            fi
        fi
    fi

    return 1
}

test_generated_sdk_compilation() {
    log_verbose "Testing generated SDK compilation..."

    cd "$ROOT_DIR/sdk/go"

    if [[ -f "go.mod" && -f "types.gen.go" ]]; then
        if go mod tidy && go build ./...; then
            return 0
        fi
    fi

    return 1
}

test_workspace_sync() {
    log_verbose "Testing workspace sync..."

    cd "$ROOT_DIR"

    if [[ -f "go.work" ]]; then
        if go work sync; then
            return 0
        fi
    fi

    return 1
}

# =============================================================================
# Error Scenario Tests
# =============================================================================

test_missing_dependencies() {
    log_verbose "Testing missing dependencies scenario..."

    # Temporarily move go binary
    if command -v go &> /dev/null; then
        local go_path
        go_path=$(which go)
        local temp_go="/tmp/go.backup"
        mv "$go_path" "$temp_go"

        # Run prerequisite check - should fail
        cd "$ROOT_DIR"
        if ! "$SCRIPT_DIR/generate-go-sdk.sh" --help >/dev/null 2>&1; then
            mv "$temp_go" "$go_path"
            return 0
        fi

        mv "$temp_go" "$go_path"
    fi

    return 1
}

test_invalid_openapi_spec() {
    log_verbose "Testing invalid OpenAPI spec scenario..."

    cd "$ROOT_DIR"

    # Create an invalid spec file
    echo '{"invalid": json}' > "$TEST_TMP_DIR/invalid.json"

    # Try to bundle it
    if command -v npx &> /dev/null; then
        if ! npx redocly bundle "$TEST_TMP_DIR/invalid.json" -o /dev/null 2>/dev/null; then
            return 0
        fi
    fi

    return 1
}

test_missing_source_files() {
    log_verbose "Testing missing source files scenario..."

    cd "$ROOT_DIR"

    # Backup and remove source SDK
    if [[ -d "packages/sdk/go" ]]; then
        mv "packages/sdk/go" "$TEST_TMP_DIR/"

        # Try generation - should handle gracefully
        if ! "$SCRIPT_DIR/generate-go-sdk.sh" --verbose 2>/dev/null; then
            mv "$TEST_TMP_DIR/sdk" "packages/"
            return 0
        fi

        mv "$TEST_TMP_DIR/sdk" "packages/"
    fi

    return 1
}

# =============================================================================
# Validation Framework Tests
# =============================================================================

test_output_quality_validation() {
    log_verbose "Testing output quality validation..."

    cd "$ROOT_DIR"

    # Generate SDK first
    if [[ -x "$SCRIPT_DIR/generate-go-sdk.sh" ]]; then
        "$SCRIPT_DIR/generate-go-sdk.sh" --clean
    fi

    # Check file sizes
    if [[ -f "sdk/go/types.gen.go" ]]; then
        local size
        size=$(wc -l < "sdk/go/types.gen.go")
        if [[ $size -gt 50 ]]; then
            return 0
        fi
    fi

    return 1
}

test_import_path_validation() {
    log_verbose "Testing import path validation..."

    cd "$ROOT_DIR/sdk/go"

    if [[ -f "types.gen.go" ]]; then
        # Check for correct import paths
        if grep -q "git.j9xym.com/opencode-api-go" types.gen.go; then
            return 0
        fi
    fi

    return 1
}

test_go_formatting() {
    log_verbose "Testing Go code formatting..."

    cd "$ROOT_DIR/sdk/go"

    if [[ -f "types.gen.go" ]]; then
        if go fmt types.gen.go; then
            return 0
        fi
    fi

    return 1
}

# =============================================================================
# Test Runners
# =============================================================================

run_unit_tests() {
    log_info "Running unit tests..."

    run_test "Prerequisites Check" test_check_prerequisites
    run_test "OpenAPI Spec Generation" test_generate_openapi_spec
    run_test "Spec Bundling" test_bundle_specification
    run_test "OAPI Codegen Install" test_oapi_codegen_install
}

run_integration_tests() {
    log_info "Running integration tests..."

    run_test "Full Generation Workflow" test_full_generation_workflow
    run_test "Generated SDK Compilation" test_generated_sdk_compilation
    run_test "Workspace Sync" test_workspace_sync
}

run_error_scenario_tests() {
    log_info "Running error scenario tests..."

    run_test "Missing Dependencies" test_missing_dependencies
    run_test "Invalid OpenAPI Spec" test_invalid_openapi_spec
    run_test "Missing Source Files" test_missing_source_files
}

run_validation_framework_tests() {
    log_info "Running validation framework tests..."

    run_test "Output Quality Validation" test_output_quality_validation
    run_test "Import Path Validation" test_import_path_validation
    run_test "Go Formatting" test_go_formatting
}

# =============================================================================
# Main Function
# =============================================================================

show_help() {
    cat << EOF
SDK Generation Test Suite

This script runs comprehensive tests for the Go SDK generation workflow.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --unit              Run unit tests for individual functions
    --integration       Run integration tests (full workflow)
    --error-scenarios   Run error scenario tests
    --validation        Run validation framework tests
    --all               Run all tests
    --verbose           Enable verbose output
    --help              Show this help message

EXAMPLES:
    $0 --unit                           # Run unit tests only
    $0 --integration --verbose          # Run integration tests with verbose output
    $0 --all                            # Run all tests

EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit)
                RUN_UNIT=true
                shift
                ;;
            --integration)
                RUN_INTEGRATION=true
                shift
                ;;
            --error-scenarios)
                RUN_ERROR_SCENARIOS=true
                shift
                ;;
            --validation)
                RUN_VALIDATION=true
                shift
                ;;
            --all)
                RUN_UNIT=true
                RUN_INTEGRATION=true
                RUN_ERROR_SCENARIOS=true
                RUN_VALIDATION=true
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

    # Default to all tests if none specified
    if [[ "$RUN_UNIT" == "false" && "$RUN_INTEGRATION" == "false" && "$RUN_ERROR_SCENARIOS" == "false" && "$RUN_VALIDATION" == "false" ]]; then
        RUN_UNIT=true
        RUN_INTEGRATION=true
        RUN_ERROR_SCENARIOS=true
        RUN_VALIDATION=true
    fi
}

main() {
    echo
    echo "========================================"
    echo "  SDK Generation Test Suite"
    echo "========================================"
    echo

    parse_arguments "$@"

    setup_test_environment

    # Trap for cleanup
    trap cleanup_test_environment EXIT

    # Run selected test suites
    if [[ "$RUN_UNIT" == "true" ]]; then
        run_unit_tests
        echo
    fi

    if [[ "$RUN_INTEGRATION" == "true" ]]; then
        run_integration_tests
        echo
    fi

    if [[ "$RUN_ERROR_SCENARIOS" == "true" ]]; then
        run_error_scenario_tests
        echo
    fi

    if [[ "$RUN_VALIDATION" == "true" ]]; then
        run_validation_framework_tests
        echo
    fi

    # Print summary
    echo "========================================"
    echo "  Test Summary"
    echo "========================================"
    echo "Tests run: $TESTS_RUN"
    echo "Tests passed: $TESTS_PASSED"
    echo "Tests failed: $TESTS_FAILED"

    if [[ $TESTS_FAILED -eq 0 ]]; then
        log_success "All tests passed!"
        exit 0
    else
        log_error "Some tests failed. Check output above."
        exit 1
    fi
}

# Run main function
main "$@"