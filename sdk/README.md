# OpenCode SDK Generation

Modern SDK generation for the OpenCode API using best-in-class tooling optimized for performance and developer experience.

---

## Overview

This directory contains **next-generation SDK generation** that replaces the previous mixed approach (hey-api + Stainless) with modern, unified tooling:

- **TypeScript SDK**: `openapi-typescript` + `openapi-fetch` (Cloudflare Workers optimized)
- **Go SDK**: `ogen` (modern Go with generics, zero-reflection)

---

## Quick start

Generate all SDKs:

```bash
cd sdk
bun install
bun run generate
```

Generate specific SDK:

```bash
bun run generate typescript  # or 'go'
```

---

## Architecture

### Current state analysis

**Modern unified SDKs:**

- ✅ `sdk/typescript/` - openapi-typescript + openapi-fetch (Cloudflare Workers optimized)
- ✅ `sdk/go/` - ogen (modern Go with generics, zero-reflection)

**Key improvements:**

- Single unified tooling approach
- Automated generation workflow
- No commercial dependencies required
- Comprehensive testing of generated SDKs
- Robust handling of complex discriminator unions

---

## New approach benefits

### TypeScript SDK advantages

**Cloudflare Workers optimized:**

- Native `fetch` API usage
- Full ESM support, zero CommonJS
- Excellent tree-shaking (import only what you need)
- Minimal bundle size impact
- Zero Node.js dependencies

**Developer experience:**

- Full TypeScript type safety from OpenAPI spec
- Auto-completion for all API endpoints
- Compile-time validation of requests/responses
- Clean, readable generated code

### Go SDK advantages

**Modern Go patterns:**

- Go 1.21+ generics support
- `Optional[T]`, `Nullable[T]` wrappers instead of pointers
- Sum types for `oneOf` with full type safety
- Zero reflection, compile-time safety
- Custom high-performance JSON codec

**Performance benefits:**

- Faster JSON parsing than standard library
- Generated validation code
- No `interface{}` usage
- Optimal HTTP client integration

---

## Directory structure

```
sdk/
├── README.md                    # This file
├── package.json                 # Root package for tooling deps
├── generate.ts                  # Main generation script
├── .gitignore                   # Ignore generated files
├── typescript/                  # TypeScript SDK
│   ├── package.json            # SDK package configuration
│   ├── generate.ts             # TS generation script
│   ├── src/
│   │   ├── index.ts           # Main export file
│   │   └── generated/         # Generated types & client
│   └── dist/                  # Built SDK (npm publishable)
└── go/                         # Go SDK
    ├── generate.ts            # Go generation script
    ├── go.mod                 # Go module
    ├── client.go              # Convenience wrapper
    ├── README.md              # Go SDK documentation
    └── generated/             # ogen generated client
```

---

## Generation workflow

### 1. OpenAPI spec generation

````bash
# From opencode server
bun run --conditions=development packages/opencode/src/index.ts generate
```</xai:function_call/>
<xai:function_call name="edit">
<parameter name="filePath">/home/jm/data/code/opencode/sdk/README.md

**Process:**

1. Generate OpenAPI spec from server
2. Install `ogen` if not present
3. Generate Go client with `ogen`
4. Create convenience wrapper and module files
5. Run `go mod tidy`

---

## Usage examples

### TypeScript (Cloudflare Workers)

```typescript
import { createOpenCodeClient } from "@opencode-ai/sdk-next"

export default {
  async fetch(request: Request): Promise<Response> {
    const client = createOpenCodeClient({
      baseUrl: "https://api.opencode.ai",
      headers: {
        Authorization: "Bearer your-token",
      },
    })

    // Fully typed API calls
    const { data, error } = await client.GET("/config")

    if (error) {
      return new Response("Error: " + error.message, { status: 500 })
    }

    return Response.json(data)
  },
}
````

### Go

```go
package main

import (
    "context"
    "log"

    opencode "github.com/sst/opencode-sdk-go"
)

func main() {
    client, err := opencode.NewClient("http://localhost:3000")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Type-safe API calls
    resp, err := client.Raw().ConfigGet(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // All responses are properly typed with generics
    fmt.Printf("Config: %+v\n", resp)
}
```

---

## Migration from legacy SDKs

The new SDKs provide a clean migration path from previous implementations. The functional API pattern ensures backward compatibility while offering enhanced type safety and performance.

### Key migration benefits:

- **Unified tooling**: Single approach for both TypeScript and Go
- **Enhanced type safety**: End-to-end validation with modern generators
- **Better performance**: Optimized for respective runtimes
- **No commercial dependencies**: Fully open source generation pipeline

---

## Development

### Add new dependencies

```bash
# TypeScript
cd sdk/typescript && bun add package-name

# Go
cd sdk/go && go get package-name
```

### Regenerate after API changes

```bash
# Regenerate all SDKs
cd sdk && bun run generate

# Or specific SDK
bun run generate typescript
```

### Testing

```bash
# TypeScript
cd sdk/typescript && bun test

# Go
cd sdk/go && go test ./...
```

---

## CI integration

**Recommended GitHub Actions workflow:**

```yaml
name: Regenerate SDKs
on:
  push:
    paths: ["packages/opencode/src/server/**"]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: oven-sh/setup-bun@v1
      - uses: actions/setup-go@v4
        with:
          go-version: "1.21"

      - run: bun install
      - run: cd sdk && bun run generate

      - name: Check for changes
        run: |
          if ! git diff --quiet; then
            echo "SDKs need regeneration"
            git add .
            git commit -m "chore: regenerate SDKs"  
            git push
          fi
```

---

## Publishing

### TypeScript SDK

```bash
cd sdk/typescript
bun run build
npm publish --access public
```

### Go SDK

Push to GitHub and tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

---

## Troubleshooting

### OpenAPI spec issues

**Discriminator problems:**

- Current spec uses complex discriminated unions
- Standard generators may fail on `discriminator` mappings
- ogen handles these better than openapi-generator

**Fix:** The new generation approach handles complex discriminators robustly.

### Generation failures

**TypeScript:**

- Ensure `openapi-typescript` CLI is available
- Check OpenAPI spec validity with `npx swagger-codegen-cli validate`

**Go:**

- Ensure Go 1.21+ installed
- Check `ogen` installation: `go install github.com/ogen-go/ogen/cmd/ogen@latest`
- Verify spec compatibility with ogen

### Performance issues

**Bundle size (TypeScript):**

- Use tree-shaking: `import { specific } from 'sdk'`
- Avoid `import * as sdk`
- Check bundle analyzer for unused code

**Memory usage (Go):**

- ogen generates efficient code by default
- Avoid holding large response objects in memory
- Use streaming for large payloads

---

## Comparison with alternatives

| Tool              | Bundle Size  | Type Safety  | CF Workers       | Maintenance    |
| ----------------- | ------------ | ------------ | ---------------- | -------------- |
| **openapi-fetch** | ✅ Tiny      | ✅ Full      | ✅ Perfect       | ✅ Active      |
| **ogen**          | ✅ Optimized | ✅ Excellent | N/A              | ✅ Active      |
| hey-api           | ⚠️ Medium    | ✅ Good      | ⚠️ Node.js deps  | ⚠️ Declining   |
| orval             | ❌ Large     | ✅ Good      | ❌ Axios-focused | ✅ Active      |
| openapi-generator | ⚠️ Verbose   | ⚠️ Limited   | N/A              | ❌ Java legacy |
| Stainless         | ✅ Good      | ✅ Good      | ✅ Good          | ✅ Commercial  |

**Recommendation:** The new approach provides the best balance of performance, type safety, and maintainability.
