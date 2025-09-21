# Opencode Go SDK

The Opencode Go library provides convenient access to the [Opencode REST API](https://opencode.ai/docs) from applications written in Go using a hybrid generation approach that combines automated type generation with manually crafted service architecture.

## Generation approach

This SDK uses a **hybrid generation approach** that balances automation with control:

- **Types**: Auto-generated using `oapi-codegen` from OpenAPI specification
- **Services**: Manual implementation following proven Stainless SDK patterns
- **Integration**: Seamless compatibility layer for existing TUI applications

---

## Requirements

- **Go 1.24+** (uses Go workspace features)
- **Node.js/Bun** (for OpenAPI spec generation)
- **oapi-codegen v2.5.0+** (automatically installed during generation)

---

## Installation

For standalone usage:

```go
go get github.com/sst/opencode-api-go
```

For development within the opencode monorepo, use the Go workspace:

```bash
# From project root
go work sync
```

---

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "github.com/sst/opencode-api-go"
)

func main() {
    client := opencode.NewClient()

    sessions, err := client.Session.List(context.TODO())
    if err != nil {
        panic(err)
    }
    fmt.Printf("Sessions: %+v\n", sessions)
}
```

---

## Development environment

### Workspace setup

The project uses Go workspaces for development. The workspace is configured in the root `go.work` file:

```
use ./sdk/go
use ./packages/tui
```

**Initialize workspace:**

```bash
go work init
go work use ./sdk/go ./packages/tui
```

**Sync workspace after changes:**

```bash
go work sync
```

---

## SDK generation workflow

### Generate from scratch

**Quick generation:**

```bash
./scripts/generate-go-sdk.sh
```

**Manual steps:**

```bash
# 1. Generate OpenAPI spec
bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json

# 2. Bundle specification
npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced

# 3. Install generation tools
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# 4. Generate types
oapi-codegen -generate types -package opencode schema/openapi.bundled.json > sdk/go/types.gen.go

# 5. Copy service layer
rsync -a --exclude 'go.mod' packages/sdk/go/ sdk/go/

# 6. Update module paths
find sdk/go -name "*.go" -exec sed -i 's/opencode.local-sdk-go/github.com/sst\/opencode-api-go/g' {} +

# 7. Tidy and build
cd sdk/go && go mod tidy && go build ./...
```

### Clean build

**Full cleanup:**

```bash
./scripts/clean-go-sdk.sh --yes
```

**Manual cleanup:**

```bash
# Remove generated artifacts
rm -f sdk/go/types.gen.go
rm -rf sdk/go/generated

# Clean module caches
go clean -modcache

# Sync workspace
go work sync
```

---

## Architecture overview

### Hybrid approach benefits

- ✅ **Type safety**: Auto-generated types from OpenAPI spec
- ✅ **Control**: Manual service layer for optimal developer experience
- ✅ **Compatibility**: Adaptor layer for existing TUI integration
- ✅ **Maintenance**: No vendor lock-in, open source toolchain

### Generated structure

```
sdk/go/
├── types.gen.go          # Auto-generated types (oapi-codegen)
├── client.go             # Main client (manual)
├── services.go           # Service registry (manual)
├── app.go                # App service (manual)
├── session.go            # Session service (manual)
├── event.go              # Event service (manual)
├── file.go               # File service (manual)
├── find.go               # Find service (manual)
├── compat/               # TUI compatibility layer
│   ├── compat.go         # Type adaptors
│   └── README.md         # Usage guide
├── internal/             # Internal utilities
└── shared/               # Shared types
```

---

## TUI integration

### Compatibility layer

The `compat` package provides adaptors for TUI integration:

```go
import "github.com/sst/opencode-api-go/compat"

// Convert generated types to TUI-compatible types
func handleConfig(generatedConfig *opencode.Config) {
    compatConfig := compat.ConvertCompatibleConfigFromGenerated(generatedConfig)
    // Use compatConfig in TUI...
}

// Convert TUI params to generated types
func logMessage(params compat.AppLogParams) error {
    generatedParams := compat.ConvertAppLogParamsToGenerated(params)
    return client.App.Log(ctx, generatedParams)
}
```

### Migration from old SDK

**Before (old SDK):**

```go
import "github.com/sst/opencode-api-go"
```

**After (hybrid SDK):**

```go
import "github.com/sst/opencode-api-go"
import "github.com/sst/opencode-api-go/compat"
```

---

## Build integration

### Generate in CI/CD

```yaml
# .github/workflows/generate-sdk.yml
- name: Generate Go SDK
  run: |
    chmod +x scripts/generate-go-sdk.sh
    ./scripts/generate-go-sdk.sh

- name: Validate SDK
  run: |
    cd sdk/go
    go mod tidy
    go build ./...
    go test ./...
```

### Validate generation

```bash
# Check generation worked
go build sdk/go/...

# Run tests
go test sdk/go/...

# Check workspace sync
go work sync
go list -m
```

---

## Differences from Stainless

| Aspect              | Hybrid SDK                  | Stainless SDK          |
| ------------------- | --------------------------- | ---------------------- |
| **Type generation** | oapi-codegen (open source)  | Stainless (commercial) |
| **Service layer**   | Manual (Stainless-inspired) | Auto-generated         |
| **Discriminators**  | Avoided (types-only)        | Full support           |
| **Customization**   | Full control                | Limited                |
| **Vendor lock-in**  | None                        | Commercial service     |
| **Maintenance**     | Self-managed                | Managed service        |

---

## Troubleshooting

### Common issues

**Build failures:**

```bash
# Clean and regenerate
./scripts/clean-go-sdk.sh --yes
./scripts/generate-go-sdk.sh

# Check Go version
go version  # Should be 1.24+
```

**Module path issues:**

```bash
# Verify correct import paths
grep -r "opencode.local-sdk-go" sdk/go/  # Should return nothing

# Fix if needed
find sdk/go -name "*.go" -exec sed -i 's/opencode.local-sdk-go/github.com/sst\/opencode-api-go/g' {} +
```

**Workspace sync problems:**

```bash
# Re-sync workspace
go work sync

# Check workspace status
go work edit -json
```

**TUI compatibility issues:**

```bash
# Check compatibility layer
cd sdk/go/compat
go build .

# Verify type conversions
go test ./compat/...
```

### Validation checklist

- [ ] `go.work` file includes both `sdk/go` and `packages/tui`
- [ ] `sdk/go/types.gen.go` exists and compiles
- [ ] Import paths use `github.com/sst/opencode-api-go`
- [ ] `go build ./...` passes from project root
- [ ] TUI application still compiles with compatibility layer

---

## Onboarding checklist

**For new developers:**

1. **Environment setup:**
   - [ ] Install Go 1.24+
   - [ ] Install Node.js/Bun
   - [ ] Clone repository
   - [ ] Run `go work sync`

2. **First generation:**
   - [ ] Run `./scripts/generate-go-sdk.sh`
   - [ ] Verify build with `go build ./...`
   - [ ] Run tests with `go test ./...`

3. **Understanding the architecture:**
   - [ ] Review `sdk/go/types.gen.go` (generated types)
   - [ ] Review `sdk/go/client.go` (manual client)
   - [ ] Review `sdk/go/compat/` (TUI integration)

4. **Development workflow:**
   - [ ] Make OpenAPI changes in TypeScript API
   - [ ] Regenerate SDK with script
   - [ ] Update compatibility layer if needed
   - [ ] Test TUI integration

---

## Best practices

### Code conventions

- **Types**: Use generated types from `types.gen.go`
- **Services**: Follow Stainless patterns for consistency
- **Errors**: Use `*opencode.Error` type for API errors
- **Context**: Always pass `context.Context` as first parameter
- **Fields**: Use `opencode.F()` helpers for request fields

### Generation workflow

- **Clean builds**: Use `clean-go-sdk.sh` before important changes
- **Validation**: Always build and test after generation
- **Backup**: Generated files are backed up automatically
- **Workspace**: Keep workspace in sync with `go work sync`

### Integration patterns

- **TUI compatibility**: Use `compat` package for type conversions
- **Error handling**: Check for `*opencode.Error` with `errors.As()`
- **Configuration**: Use workspace for development, modules for production
