#!/bin/bash
# =============================================================================
# SDK Quality Validation Framework
# =============================================================================
#
# Validates the quality and correctness of generated Go SDK
# Checks for compilation, formatting, security, and performance issues
#
# Usage:
#   ./scripts/validate-sdk-quality.sh [options]
#
# Options:
#   --build          Validate build compilation
#   --format         Check code formatting
#   --security       Security vulnerability checks
#   --performance    Performance validation
#   --quality        Code quality metrics
#   --all            Run all validations
#   --verbose        Enable verbose output
#   --help           Show this help
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

# Validation flags
CHECK_BUILD=false
CHECK_FORMAT=false
CHECK_SECURITY=false
CHECK_PERFORMANCE=false
CHECK_QUALITY=false
VERBOSE=false

# Validation results
ISSUES_FOUND=0
ISSUES_CRITICAL=0
ISSUES_WARNING=0

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

report_issue() {
    local level="$1"
    local message="$2"
    local details="${3:-}"

    ((ISSUES_FOUND++))

    case "$level" in
        "CRITICAL")
            ((ISSUES_CRITICAL++))
            log_error "CRITICAL: $message"
            ;;
        "WARNING")
            ((ISSUES_WARNING++))
            log_warning "WARNING: $message"
            ;;
        "INFO")
            log_info "INFO: $message"
            ;;
    esac

    if [[ -n "$details" ]]; then
        echo "  $details"
    fi
}

# =============================================================================
# Build Validation
# =============================================================================

validate_build() {
    log_info "Validating SDK build..."

    if [[ ! -d "$SDK_DIR" ]]; then
        report_issue "CRITICAL" "SDK directory not found" "$SDK_DIR"
        return 1
    fi

    cd "$SDK_DIR"

    # Check for required files
    local required_files=("go.mod" "types.gen.go")
    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            report_issue "CRITICAL" "Required file missing" "$file"
        fi
    done

    # Validate go.mod
    if [[ -f "go.mod" ]]; then
        if ! go mod tidy; then
            report_issue "CRITICAL" "go mod tidy failed"
        fi

        # Check for module name
        if ! grep -q "^module " go.mod; then
            report_issue "CRITICAL" "Module declaration missing in go.mod"
        fi
    fi

    # Attempt to build
    if ! go build ./...; then
        report_issue "CRITICAL" "SDK build failed"

        # Analyze build errors
        local build_errors
        build_errors=$(go build ./... 2>&1 || true)

        if echo "$build_errors" | grep -q "redeclared in this block"; then
            report_issue "CRITICAL" "Type redeclaration errors found"
        fi

        if echo "$build_errors" | grep -q "undefined:"; then
            report_issue "CRITICAL" "Undefined symbol errors found"
        fi

        if echo "$build_errors" | grep -q "cannot find package"; then
            report_issue "WARNING" "Missing package dependencies"
        fi

        if [[ "$VERBOSE" == "true" ]]; then
            echo "Build errors:"
            echo "$build_errors"
        fi
    else
        log_success "SDK builds successfully"
    fi

    # Run go vet
    if ! go vet ./...; then
        report_issue "WARNING" "go vet found issues"
    fi

    # Check for unused imports
    local unused_imports
    unused_imports=$(go mod tidy 2>&1 | grep -c "unused" || echo "0")
    if [[ $unused_imports -gt 0 ]]; then
        report_issue "WARNING" "Unused imports detected" "$unused_imports unused imports"
    fi
}

# =============================================================================
# Code Formatting Validation
# =============================================================================

validate_formatting() {
    log_info "Validating code formatting..."

    cd "$SDK_DIR"

    # Check gofmt
    local fmt_issues
    fmt_issues=$(gofmt -l . | wc -l)

    if [[ $fmt_issues -gt 0 ]]; then
        report_issue "WARNING" "Code formatting issues found" "$fmt_issues files need formatting"
        if [[ "$VERBOSE" == "true" ]]; then
            log_verbose "Files needing formatting:"
            gofmt -l . | sed 's/^/  /'
        fi
    else
        log_success "Code is properly formatted"
    fi

    # Check goimports if available
    if command -v goimports &> /dev/null; then
        local imports_issues
        imports_issues=$(goimports -l . | wc -l)

        if [[ $imports_issues -gt 0 ]]; then
            report_issue "WARNING" "Import formatting issues found" "$imports_issues files need import formatting"
        fi
    fi

    # Check for long lines
    local long_lines
    long_lines=$(find . -name "*.go" -exec grep -l '.\{121,\}' {} \; | wc -l)

    if [[ $long_lines -gt 0 ]]; then
        report_issue "INFO" "Long lines found" "$long_lines files have lines longer than 120 characters"
    fi
}

# =============================================================================
# Security Validation
# =============================================================================

validate_security() {
    log_info "Validating security..."

    cd "$SDK_DIR"

    # Check for hardcoded secrets
    local secret_patterns=("password" "secret" "token" "key" "apikey" "auth")
    local secret_found=false

    for pattern in "${secret_patterns[@]}"; do
        if grep -r -i "$pattern" --include="*.go" . | grep -v -E "(test|example|comment|doc)" | grep -v "import"; then
            secret_found=true
            break
        fi
    done

    if [[ "$secret_found" == "true" ]]; then
        report_issue "CRITICAL" "Potential hardcoded secrets found"
    fi

    # Check for insecure practices
    if grep -r "http://" --include="*.go" . | grep -v "example\|test\|comment"; then
        report_issue "WARNING" "HTTP URLs found (should use HTTPS)"
    fi

    # Check for unsafe pointer usage
    if grep -r "unsafe\." --include="*.go" .; then
        report_issue "INFO" "Unsafe package usage found"
    fi

    # Check for SQL injection vulnerabilities (basic check)
    if grep -r -i "exec\|query" --include="*.go" . | grep -v "test\|example"; then
        report_issue "INFO" "Database operations found - review for SQL injection"
    fi

    # Check for proper error handling
    local error_checks
    error_checks=$(grep -r "if err != nil" --include="*.go" . | wc -l)
    local total_functions
    total_functions=$(grep -r "^func " --include="*.go" . | wc -l)

    if [[ $total_functions -gt 0 ]]; then
        local error_ratio=$((error_checks * 100 / total_functions))
        if [[ $error_ratio -lt 50 ]]; then
            report_issue "INFO" "Low error handling coverage" "$error_checks error checks in $total_functions functions"
        fi
    fi
}

# =============================================================================
# Performance Validation
# =============================================================================

validate_performance() {
    log_info "Validating performance aspects..."

    cd "$SDK_DIR"

    # Check for inefficient operations
    if grep -r "append.*append" --include="*.go" .; then
        report_issue "INFO" "Potential inefficient append operations"
    fi

    # Check for string concatenation in loops
    if grep -A5 -B5 "for.*range" --include="*.go" . | grep -E "[+]=.*string"; then
        report_issue "INFO" "String concatenation in loops detected"
    fi

    # Check for large struct allocations
    local large_structs
    large_structs=$(grep -A 20 "type.*struct" --include="*.go" . | grep -c "string\|int\|bool" || echo "0")

    if [[ $large_structs -gt 50 ]]; then
        report_issue "INFO" "Large structs detected - consider optimization"
    fi

    # Check for memory leaks (basic heuristic)
    if grep -r "new(\|make(" --include="*.go" . | grep -v "test\|example"; then
        report_issue "INFO" "Memory allocations found - review for potential leaks"
    fi

    # Check bundle size
    if [[ -f "types.gen.go" ]]; then
        local file_size
        file_size=$(wc -l < "types.gen.go")

        if [[ $file_size -gt 5000 ]]; then
            report_issue "INFO" "Large generated file" "types.gen.go has $file_size lines - consider optimization"
        fi
    fi
}

# =============================================================================
# Code Quality Validation
# =============================================================================

validate_quality() {
    log_info "Validating code quality..."

    cd "$SDK_DIR"

    # Check for generated file markers
    if [[ -f "types.gen.go" ]]; then
        if ! grep -q "Code generated" "types.gen.go"; then
            report_issue "WARNING" "Generated file missing header comment"
        fi
    fi

    # Check for documentation
    local undocumented_exports
    undocumented_exports=$(grep -r "^type\|^func" --include="*.go" . | grep -v "^//" | wc -l)

    if [[ $undocumented_exports -gt 10 ]]; then
        report_issue "INFO" "Many undocumented exported types/functions" "$undocumented_exports exported items without documentation"
    fi

    # Check for TODO/FIXME comments
    local todos
    todos=$(grep -r -i "todo\|fixme\|hack" --include="*.go" . | wc -l)

    if [[ $todos -gt 0 ]]; then
        report_issue "INFO" "TODO/FIXME comments found" "$todos items need attention"
    fi

    # Check for proper naming conventions
    if grep -r "^[a-z].*[A-Z]" --include="*.go" . | grep -E "^func|^type|^var"; then
        report_issue "INFO" "Potential naming convention issues"
    fi

    # Check for magic numbers
    local magic_numbers
    magic_numbers=$(grep -r -P "\b\d{2,}\b" --include="*.go" . | grep -v -E "(const|var|return|if|for|case)" | wc -l)

    if [[ $magic_numbers -gt 20 ]]; then
        report_issue "INFO" "Magic numbers detected" "$magic_numbers potential magic numbers"
    fi

    # Check for proper test coverage (if tests exist)
    if ls *_test.go >/dev/null 2>&1; then
        local test_coverage
        if go test -cover ./... >/dev/null 2>&1; then
            test_coverage=$(go test -cover ./... 2>/dev/null | grep -o "coverage: [0-9.]*%" | grep -o "[0-9.]*" | head -1 || echo "0")
            if [[ $(echo "$test_coverage < 50" | bc -l) -eq 1 ]]; then
                report_issue "INFO" "Low test coverage" "${test_coverage}% coverage"
            fi
        fi
    else
        report_issue "INFO" "No test files found"
    fi
}

# =============================================================================
# Main Function
# =============================================================================

show_help() {
    cat << EOF
SDK Quality Validation Framework

This script validates the quality and correctness of the generated Go SDK.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --build          Validate build compilation
    --format         Check code formatting
    --security       Security vulnerability checks
    --performance    Performance validation
    --quality        Code quality metrics
    --all            Run all validations
    --verbose        Enable verbose output
    --help           Show this help message

EXAMPLES:
    $0 --build --format                    # Check build and formatting
    $0 --all --verbose                     # Run all validations with verbose output

VALIDATION RESULTS:
    - CRITICAL: Issues that prevent the SDK from working
    - WARNING: Issues that should be addressed
    - INFO: Suggestions for improvement

EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --build)
                CHECK_BUILD=true
                shift
                ;;
            --format)
                CHECK_FORMAT=true
                shift
                ;;
            --security)
                CHECK_SECURITY=true
                shift
                ;;
            --performance)
                CHECK_PERFORMANCE=true
                shift
                ;;
            --quality)
                CHECK_QUALITY=true
                shift
                ;;
            --all)
                CHECK_BUILD=true
                CHECK_FORMAT=true
                CHECK_SECURITY=true
                CHECK_PERFORMANCE=true
                CHECK_QUALITY=true
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

    # Default to all checks if none specified
    if [[ "$CHECK_BUILD" == "false" && "$CHECK_FORMAT" == "false" && "$CHECK_SECURITY" == "false" && "$CHECK_PERFORMANCE" == "false" && "$CHECK_QUALITY" == "false" ]]; then
        CHECK_BUILD=true
        CHECK_FORMAT=true
        CHECK_SECURITY=true
        CHECK_PERFORMANCE=true
        CHECK_QUALITY=true
    fi
}

main() {
    echo
    echo "========================================"
    echo "  SDK Quality Validation Framework"
    echo "========================================"
    echo

    parse_arguments "$@"

    if [[ ! -d "$SDK_DIR" ]]; then
        log_error "SDK directory not found: $SDK_DIR"
        log_info "Run SDK generation first: ./scripts/generate-go-sdk.sh"
        exit 1
    fi

    # Run selected validations
    if [[ "$CHECK_BUILD" == "true" ]]; then
        validate_build
        echo
    fi

    if [[ "$CHECK_FORMAT" == "true" ]]; then
        validate_formatting
        echo
    fi

    if [[ "$CHECK_SECURITY" == "true" ]]; then
        validate_security
        echo
    fi

    if [[ "$CHECK_PERFORMANCE" == "true" ]]; then
        validate_performance
        echo
    fi

    if [[ "$CHECK_QUALITY" == "true" ]]; then
        validate_quality
        echo
    fi

    # Print summary
    echo "========================================"
    echo "  Validation Summary"
    echo "========================================"
    echo "Total issues found: $ISSUES_FOUND"
    echo "Critical issues: $ISSUES_CRITICAL"
    echo "Warning issues: $ISSUES_WARNING"
    echo "Info issues: $((ISSUES_FOUND - ISSUES_CRITICAL - ISSUES_WARNING))"

    if [[ $ISSUES_CRITICAL -gt 0 ]]; then
        log_error "Critical issues found - SDK may not work correctly"
        exit 1
    elif [[ $ISSUES_WARNING -gt 0 ]]; then
        log_warning "Warning issues found - consider addressing them"
        exit 0
    else
        log_success "No critical issues found!"
        exit 0
    fi
}

# Run main function
main "$@"