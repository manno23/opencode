#!/bin/bash
# =============================================================================
# Docker-based SDK Testing Script
# =============================================================================
#
# Runs comprehensive SDK tests in Docker containers for consistent environments
#
# Usage:
#   ./scripts/run-docker-sdk-tests.sh [options]
#
# Options:
#   --all           Run all Docker test services
#   --linux         Run Linux environment tests
#   --clean         Run clean environment tests
#   --multiarch     Run multi-architecture tests
#   --validation    Run validation-only tests
#   --build         Build Docker images only
#   --cleanup       Clean up Docker containers and images
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
readonly COMPOSE_FILE="$ROOT_DIR/docker-compose.sdk-test.yml"

# Determine container and compose commands
CONTAINER_CMD=""
COMPOSE_CMD=""

check_container_tool() {
    if command -v docker &> /dev/null; then
        CONTAINER_CMD="docker"
    elif command -v podman &> /dev/null; then
        CONTAINER_CMD="podman"
    else
        log_error "Neither Docker nor Podman is installed or available"
        exit 1
    fi

    if [[ "$CONTAINER_CMD" == "docker" ]]; then
        if command -v docker-compose &> /dev/null; then
            COMPOSE_CMD="docker-compose"
        else
            log_error "docker-compose is required when using Docker"
            exit 1
        fi
    else
        if command -v podman-compose &> /dev/null; then
            COMPOSE_CMD="podman-compose"
        else
            COMPOSE_CMD="podman compose"
        fi
    fi
}

# Test flags
RUN_ALL=false
RUN_LINUX=false
RUN_CLEAN=false
RUN_MULTIARCH=false
RUN_VALIDATION=false
BUILD_ONLY=false
CLEANUP=false
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



cleanup_containers() {
    log_info "Cleaning up Docker containers and images..."

    # Stop and remove containers
    $COMPOSE_CMD -f "$COMPOSE_FILE" down --remove-orphans || true

    # Remove images
    $COMPOSE_CMD -f "$COMPOSE_FILE" rm -f || true

    # Clean up dangling images
    $CONTAINER_CMD image prune -f || true

    log_success "Cleanup completed"
}

build_images() {
    log_info "Building Docker images..."

    if ! $COMPOSE_CMD -f "$COMPOSE_FILE" build; then
        log_error "Failed to build Docker images"
        exit 1
    fi

    log_success "Docker images built successfully"
}

run_service() {
    local service="$1"
    local description="$2"

    log_info "Running $description..."

    if [[ "$VERBOSE" == "true" ]]; then
        $COMPOSE_CMD -f "$COMPOSE_FILE" up "$service"
    else
        $COMPOSE_CMD -f "$COMPOSE_FILE" up --abort-on-container-exit "$service"
    fi

    local exit_code=$?
    if [[ $exit_code -eq 0 ]]; then
        log_success "$description completed successfully"
    else
        log_error "$description failed with exit code $exit_code"
        return $exit_code
    fi
}

# =============================================================================
# Test Functions
# =============================================================================

run_linux_tests() {
    run_service "sdk-test-linux" "Linux SDK tests"
}

run_clean_tests() {
    run_service "sdk-test-clean-env" "Clean environment tests"
}

run_multiarch_tests() {
    run_service "sdk-test-multiarch" "Multi-architecture tests"
}

run_validation_tests() {
    run_service "sdk-validation-only" "SDK validation tests"
}

run_all_tests() {
    log_info "Running all Docker-based SDK tests..."

    local failed_services=()

    # Run each service and track failures
    if ! run_service "sdk-test-linux" "Linux SDK tests"; then
        failed_services+=("linux")
    fi

    if ! run_service "sdk-test-clean-env" "Clean environment tests"; then
        failed_services+=("clean")
    fi

    if ! run_service "sdk-test-multiarch" "Multi-architecture tests"; then
        failed_services+=("multiarch")
    fi

    if ! run_service "sdk-validation-only" "SDK validation tests"; then
        failed_services+=("validation")
    fi

    # Report results
    if [[ ${#failed_services[@]} -eq 0 ]]; then
        log_success "All Docker tests passed!"
    else
        log_error "Some tests failed: ${failed_services[*]}"
        return 1
    fi
}

# =============================================================================
# Main Function
# =============================================================================

show_help() {
    cat << EOF
Docker-based SDK Testing Script

This script runs comprehensive SDK tests in Docker containers for consistent environments.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --all           Run all Docker test services
    --linux         Run Linux environment tests
    --clean         Run clean environment tests
    --multiarch     Run multi-architecture tests
    --validation    Run validation-only tests
    --build         Build Docker images only
    --cleanup       Clean up Docker containers and images
    --verbose       Enable verbose output
    --help          Show this help message

SERVICES:
    sdk-test-linux      - Full test suite in Linux environment
    sdk-test-clean-env  - Clean environment with generation and validation
    sdk-test-multiarch  - Multi-architecture compatibility tests
    sdk-validation-only - SDK quality validation only

EXAMPLES:
    $0 --all                      # Run all tests
    $0 --linux --verbose          # Run Linux tests with verbose output
    $0 --build                    # Build images only
    $0 --cleanup                  # Clean up containers and images

EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --all)
                RUN_ALL=true
                shift
                ;;
            --linux)
                RUN_LINUX=true
                shift
                ;;
            --clean)
                RUN_CLEAN=true
                shift
                ;;
            --multiarch)
                RUN_MULTIARCH=true
                shift
                ;;
            --validation)
                RUN_VALIDATION=true
                shift
                ;;
            --build)
                BUILD_ONLY=true
                shift
                ;;
            --cleanup)
                CLEANUP=true
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
    if [[ "$RUN_ALL" == "false" && "$RUN_LINUX" == "false" && "$RUN_CLEAN" == "false" && "$RUN_MULTIARCH" == "false" && "$RUN_VALIDATION" == "false" && "$BUILD_ONLY" == "false" && "$CLEANUP" == "false" ]]; then
        RUN_ALL=true
    fi
}

main() {
    echo
    echo "========================================"
    echo "  Docker SDK Testing Environment"
    echo "========================================"
    echo

    parse_arguments "$@"

    check_container_tool

    cd "$ROOT_DIR"

    # Handle cleanup first
    if [[ "$CLEANUP" == "true" ]]; then
        cleanup_containers
        exit 0
    fi

    # Build images if needed
    if [[ "$BUILD_ONLY" == "true" ]]; then
        build_images
        exit 0
    fi

    # Run selected tests
    local exit_code=0

    if [[ "$RUN_ALL" == "true" ]]; then
        run_all_tests || exit_code=$?
    fi

    if [[ "$RUN_LINUX" == "true" ]]; then
        run_linux_tests || exit_code=$?
    fi

    if [[ "$RUN_CLEAN" == "true" ]]; then
        run_clean_tests || exit_code=$?
    fi

    if [[ "$RUN_MULTIARCH" == "true" ]]; then
        run_multiarch_tests || exit_code=$?
    fi

    if [[ "$RUN_VALIDATION" == "true" ]]; then
        run_validation_tests || exit_code=$?
    fi

    # Final cleanup
    if [[ $exit_code -eq 0 ]]; then
        log_success "All requested tests completed successfully!"
    else
        log_error "Some tests failed. Check output above."
    fi

    exit $exit_code
}

# Run main function
main "$@"