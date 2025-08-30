#!/bin/bash
# =============================================================================
# Clean SDK Generation Script
# =============================================================================
#
# Generate SDKs from scratch without relying on packages/sdk/js
# This creates a completely independent SDK generation process
#
# Usage:
#   ./scripts/generate-clean-sdk.sh [options]
#
# Options:
#   --typescript    Generate TypeScript SDK only
#   --go           Generate Go SDK only
#   --all          Generate both SDKs (default)
#   --server       Start server for generation
#   --help         Show this help
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
readonly SDK_DIR="$ROOT_DIR/sdk"
readonly SERVER_PID_FILE="$ROOT_DIR/.server.pid"

# Flags
GENERATE_TS=false
GENERATE_GO=false
START_SERVER=false

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

cleanup_server() {
    if [[ -f "$SERVER_PID_FILE" ]]; then
        local pid
        pid=$(cat "$SERVER_PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            log_info "Stopping server (PID: $pid)"
            kill "$pid"
            sleep 2
        fi
        rm -f "$SERVER_PID_FILE"
    fi
}

start_server() {
    log_info "Starting OpenCode server for SDK generation..."

    # Start server in background
    cd "$ROOT_DIR"
    opencode serve --port 4096 --hostname 127.0.0.1 > /dev/null 2>&1 &
    echo $! > "$SERVER_PID_FILE"

    # Wait for server to be ready
    local retries=30
    while [[ $retries -gt 0 ]]; do
        if curl -s http://localhost:4096/app > /dev/null 2>&1; then
            log_success "Server is ready"
            return 0
        fi
        sleep 1
        ((retries--))
    done

    log_error "Server failed to start"
    cleanup_server
    return 1
}

generate_openapi_spec() {
    log_info "Generating OpenAPI spec from server..."

    cd "$ROOT_DIR"

    # Generate spec from running server
    curl -s http://localhost:4096/doc > schema/openapi.json

    # Validate the spec
    if [[ ! -s "schema/openapi.json" ]]; then
        log_error "Failed to generate OpenAPI spec"
        return 1
    fi

    log_success "OpenAPI spec generated"
}

generate_typescript_sdk() {
    log_info "Generating TypeScript SDK..."

    cd "$ROOT_DIR"

    # Create SDK directory
    mkdir -p "$SDK_DIR/typescript/src"

    # Generate TypeScript types
    npx openapi-typescript schema/openapi.json --output "$SDK_DIR/typescript/src/types.ts"

    # Generate client
    cat > "$SDK_DIR/typescript/src/client.ts" << 'EOF'
import createClient from 'openapi-fetch'
import type { paths } from './types.js'

export type OpenCodeClient = ReturnType<typeof createClient<paths>>

export function createOpenCodeClient(options: {
  baseUrl?: string
  headers?: Record<string, string>
  fetch?: typeof fetch
} = {}): OpenCodeClient {
  const {
    baseUrl = 'http://localhost:4096',
    headers = {},
    fetch: customFetch
  } = options

  return createClient<paths>({
    baseUrl,
    headers: {
      'User-Agent': 'opencode-sdk-typescript/1.0.0',
      ...headers
    },
    ...(customFetch && { fetch: customFetch })
  })
}

export * from './types.js'
EOF

    # Generate package.json
    cat > "$SDK_DIR/typescript/package.json" << 'EOF'
{
  "name": "@opencode-ai/sdk-local",
  "version": "1.0.0",
  "type": "module",
  "main": "./src/client.ts",
  "types": "./src/types.ts",
  "exports": {
    ".": {
      "types": "./src/types.ts",
      "import": "./src/client.ts"
    }
  },
  "dependencies": {
    "openapi-fetch": "^0.12.0"
  },
  "devDependencies": {
    "openapi-typescript": "^7.5.0",
    "typescript": "^5.8.2"
  }
}
EOF

    # Generate tsconfig.json
    cat > "$SDK_DIR/typescript/tsconfig.json" << 'EOF'
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "declaration": true,
    "declarationMap": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": false,
    "outDir": "dist"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
EOF

    log_success "TypeScript SDK generated in $SDK_DIR/typescript/"
}

generate_go_sdk() {
    log_info "Generating Go SDK..."

    cd "$ROOT_DIR"

    # Bundle the OpenAPI spec
    npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json

    # Generate Go types
    go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
        -generate types \
        -package opencode \
        schema/openapi.bundled.json > "$SDK_DIR/go/types.gen.go"

    # Create basic Go client
    cat > "$SDK_DIR/go/client.go" << 'EOF'
package opencode

import (
    "context"
    "net/http"
)

type Client struct {
    server string
    httpClient *http.Client
}

func NewClient(server string) *Client {
    return &Client{
        server: server,
        httpClient: &http.Client{},
    }
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    // Basic HTTP client implementation
    // This would be expanded based on the API endpoints
    return nil, nil
}
EOF

    # Create go.mod
    cat > "$SDK_DIR/go/go.mod" << 'EOF'
module git.j9xym.com/opencode-api-go

go 1.24

require (
    github.com/oapi-codegen/oapi-codegen/v2 v2.4.1
)
EOF

    # Tidy Go modules
    cd "$SDK_DIR/go"
    go mod tidy
    go fmt ./...
    cd "$ROOT_DIR"

    log_success "Go SDK generated in $SDK_DIR/go/"
}

show_help() {
    cat << EOF
Clean SDK Generation Script

Generate SDKs from scratch without relying on packages/sdk/js.
This creates a completely independent SDK generation process.

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --typescript    Generate TypeScript SDK only
    --go           Generate Go SDK only
    --all          Generate both SDKs (default)
    --server       Start server automatically for generation
    --help         Show this help message

EXAMPLES:
    $0                          # Generate both SDKs
    $0 --typescript            # Generate only TypeScript SDK
    $0 --go --server           # Start server and generate only Go SDK
    $0 --all --server          # Start server and generate both SDKs

OUTPUT:
    TypeScript SDK: sdk/typescript/
    Go SDK:         sdk/go/

This script is completely independent of packages/sdk/js/
EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --typescript)
                GENERATE_TS=true
                shift
                ;;
            --go)
                GENERATE_GO=true
                shift
                ;;
            --all)
                GENERATE_TS=true
                GENERATE_GO=true
                shift
                ;;
            --server)
                START_SERVER=true
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

    # Default to generating both if nothing specified
    if [[ "$GENERATE_TS" == "false" && "$GENERATE_GO" == "false" ]]; then
        GENERATE_TS=true
        GENERATE_GO=true
    fi
}

main() {
    echo
    echo "========================================"
    echo "  🧹 Clean SDK Generation"
    echo "========================================"
    echo
    log_info "Generating SDKs independently of packages/sdk/js/"
    echo

    parse_arguments "$@"

    # Trap for cleanup
    trap cleanup_server EXIT

    # Start server if requested
    if [[ "$START_SERVER" == "true" ]]; then
        start_server
    fi

    # Generate OpenAPI spec
    generate_openapi_spec

    # Generate TypeScript SDK
    if [[ "$GENERATE_TS" == "true" ]]; then
        generate_typescript_sdk
    fi

    # Generate Go SDK
    if [[ "$GENERATE_GO" == "true" ]]; then
        generate_go_sdk
    fi

    echo
    log_success "Clean SDK generation complete!"
    echo
    echo "Generated SDKs (completely independent of packages/sdk/js/):"
    if [[ "$GENERATE_TS" == "true" ]]; then
        echo "• TypeScript: $SDK_DIR/typescript/"
    fi
    if [[ "$GENERATE_GO" == "true" ]]; then
        echo "• Go:         $SDK_DIR/go/"
    fi
    echo
    echo "These SDKs are generated fresh from your server API! 🎉"
}

# Run main function
main "$@"