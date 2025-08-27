#!/bin/bash
# Generate Go SDK from OpenAPI spec and prepare service layer
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"

echo "Generating OpenAPI spec..."
bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json

echo "Bundling spec..."
npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced

echo "Installing oapi-codegen..."
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

echo "Generating Go types..."
oapi-codegen -generate types -package opencode schema/openapi.bundled.json > sdk/go/types.gen.go

echo "Copying service layer from packages/sdk/go..."
rsync -a --exclude 'go.mod' packages/sdk/go/ sdk/go/

echo "Updating import paths to module git.j9xym.com/opencode-api-go..."
rg --vimgrep "opencode.local-sdk-go" -g "sdk/go/**" || true
find sdk/go -type f -name "*.go" -exec sed -i 's/opencode.local-sdk-go/git.j9xym.com\/opencode-api-go/g' {} +

echo "Tidying modules..."
cd sdk/go && go mod tidy

echo "Build SDK..."
go build ./...

echo "Sync workspace if go.work exists..."
cd "$ROOT"
if [ -f go.work ]; then
  go work sync
fi

echo "Generation complete"
