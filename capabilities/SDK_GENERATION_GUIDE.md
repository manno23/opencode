# OpenAPI SDK Generation Guide

This document provides a comprehensive guide to generating SDKs from OpenAPI specifications for the opencode project, focusing on the recommended hybrid approach that combines the best of automated generation with proven SDK architecture patterns.

## Overview

The opencode project uses a **hybrid approach** for SDK generation:

1. **TypeScript**: `@hey-api/openapi-ts` for modern, type-safe SDKs
2. **Go**: `oapi-codegen` for type generation + manual Stainless-inspired service architecture

## OpenAPI Specification

### Generated Specification

**Location**: `schema/openapi.json`  
**Generation Command**:

```bash
bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json
```

**Purpose**: Code-first OpenAPI specification generated from the TypeScript API implementation.

### Bundled Specification (For Generation)

**Location**: `schema/openapi.bundled.json`  
**Generation Command**:

```bash
npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced
```

**Purpose**: Dereferenced specification with all `$ref` references resolved inline - required for SDK generation tools.

## TypeScript SDK Generation

### Recommended Approach: @hey-api/openapi-ts

**Location**: `sdk/typescript/`

**Generation Command**:

```bash
npx @hey-api/openapi-ts --input schema/openapi.bundled.json --output sdk/typescript --client fetch
```

**Configuration**:

```json
{
  "dependencies": {
    "@hey-api/client-fetch": "^0.1.0"
  }
}
```

**Generated Structure**:

```
sdk/typescript/
├── types.gen.ts     # TypeScript type definitions
├── sdk.gen.ts       # SDK method implementations
├── client.gen.ts    # HTTP client
└── index.ts         # Main exports
```

**Features**:

- ✅ Modern TypeScript with full type safety
- ✅ Handles discriminated unions correctly
- ✅ Clean, idiomatic generated code
- ✅ Browser and Node.js support

## Go SDK Generation

### Recommended Approach: Hybrid (oapi-codegen + Manual Services)

**Location**: `sdk/go/`

**Step 1: Generate Types with oapi-codegen**

```bash
# Install oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Generate clean Go types
oapi-codegen -generate types -package opencode schema/openapi.bundled.json > sdk/go/types.gen.go
```

**Step 2: Manual Service Structure (Based on Stainless Patterns)**

```go
// client.go
type Client struct {
    Options []RequestOption
    App     *AppService
    Session *SessionService
    Event   *EventService
    File    *FileService
    Find    *FindService
}
```

**Configuration File**: `sdk/go/oapi-config.yaml`

```yaml
package: opencode
generate: [types]
output: types.gen.go
```

**Why This Approach Works**:

- ✅ **oapi-codegen types**: Clean, idiomatic Go types without discriminator issues
- ✅ **Manual services**: Full control over API structure and developer experience
- ✅ **Stainless patterns**: Copy proven service architecture from `packages/sdk/go/`
- ✅ **No discriminator problems**: Types-only generation avoids complex discriminator edge cases

## Approach Comparison

| Aspect                    | Hybrid (Recommended)         | @hey-api/openapi-ts | Legacy Stainless      |
| ------------------------- | ---------------------------- | ------------------- | --------------------- |
| **Go Support**            | ✅ oapi-codegen types        | ❌ TypeScript only  | ✅ Full service SDK   |
| **TypeScript Support**    | ⚠️ Manual                    | ✅ Excellent        | ⚠️ Limited            |
| **Discriminator Support** | ✅ Types only (avoids issue) | ✅ Excellent        | ✅ Excellent          |
| **Code Quality**          | ✅ Full control              | ✅ Production-ready | ✅ Production-ready   |
| **Maintenance**           | ✅ Self-managed              | ✅ Open source      | 💰 Commercial service |
| **Setup Complexity**      | ⚠️ Manual service layer      | ✅ Simple           | ⚠️ Requires account   |
| **Build Integration**     | ✅ Scriptable                | ✅ Scriptable       | ✅ CI/CD friendly     |

## Practical Generation Scripts

### Complete TypeScript SDK Generation

```bash
#!/bin/bash
# generate-typescript-sdk.sh

# 1. Generate OpenAPI spec
echo "Generating OpenAPI specification..."
bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json

# 2. Bundle specification
echo "Bundling specification..."
npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced

# 3. Generate SDK
echo "Generating TypeScript SDK..."
npx @hey-api/openapi-ts --input schema/openapi.bundled.json --output sdk/typescript --client fetch

# 4. Install dependencies
echo "Installing SDK dependencies..."
cd sdk/typescript && npm install

echo "TypeScript SDK generation complete!"
```

### Complete Go SDK Generation

```bash
#!/bin/bash
# generate-go-sdk.sh

# 1. Generate OpenAPI spec (reuse from TypeScript generation)
if [ ! -f "schema/openapi.bundled.json" ]; then
    echo "Generating OpenAPI specification..."
    bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json
    npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced
fi

# 2. Generate types with oapi-codegen
echo "Generating Go types..."
oapi-codegen -generate types -package opencode schema/openapi.bundled.json > sdk/go/types.gen.go

# 3. Copy service structure from Stainless SDK (packages/sdk/go/)
echo "Setting up service structure..."
cp -r packages/sdk/go/{client.go,app.go,session.go,event.go,file.go,find.go} sdk/go/
cp -r packages/sdk/go/{option,internal,packages} sdk/go/

# 4. Update imports
echo "Updating import paths..."
find sdk/go -name "*.go" -exec sed -i 's/opencode\.local-sdk-go/github.com\/sst\/opencode-sdk-go/g' {} \;

# 5. Update go.mod
echo "Updating Go module..."
cd sdk/go
go mod init github.com/sst/opencode-sdk-go
go mod tidy

echo "Go SDK generation complete!"
```

### Configuration Files

**TypeScript Config**: `sdk/typescript/tsconfig.json`

```json
{
  "extends": "@tsconfig/node22/tsconfig.json",
  "compilerOptions": {
    "outDir": "./dist",
    "rootDir": "./src"
  },
  "include": ["**/*.ts"],
  "exclude": ["node_modules", "dist"]
}
```

**Go Config**: `sdk/go/oapi-config.yaml`

```yaml
package: opencode
generate: [types]
output: types.gen.go
```

## Quick Start

### Generate TypeScript SDK

```bash
chmod +x scripts/generate-typescript-sdk.sh
./scripts/generate-typescript-sdk.sh
```

### Generate Go SDK

```bash
chmod +x scripts/generate-go-sdk.sh
./scripts/generate-go-sdk.sh
```

## Recommendations

**Production Use:**

- **TypeScript**: `@hey-api/openapi-ts` - Modern, type-safe, handles discriminators
- **Go**: Hybrid approach - oapi-codegen types + manual Stainless-inspired services

**Key Advantages:**

- ✅ **No discriminator issues** - Types-only generation avoids complex edge cases
- ✅ **Full control** - Manual service layer provides excellent developer experience
- ✅ **Open source** - No commercial dependencies or vendor lock-in
- ✅ **Battle-tested patterns** - Leverage proven Stainless SDK architecture

## Troubleshooting

### Common Issues

**TypeScript Generation:**

```bash
# Missing dependencies
npm install @hey-api/client-fetch

# Bundling issues
npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced
```

**Go Generation:**

```bash
# Install oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Missing types
oapi-codegen -generate types -package opencode schema/openapi.bundled.json > sdk/go/types.gen.go
```

**Import Path Issues:**

```bash
# Fix import paths in copied service files
find sdk/go -name "*.go" -exec sed -i 's/opencode\.local-sdk-go/github.com\/sst\/opencode-sdk-go/g' {} \;

# Update go.mod
cd sdk/go && go mod init github.com/sst/opencode-sdk-go && go mod tidy
```

### Validation

```bash
# Validate OpenAPI spec
npx redocly lint schema/openapi.json

# Test TypeScript compilation
cd sdk/typescript && npx tsc --noEmit

# Test Go compilation
cd sdk/go && go build ./...
```
