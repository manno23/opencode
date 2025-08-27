# OpenCode Go SDK

Modern Go SDK for the OpenCode API, generated using [ogen](https://github.com/ogen-go/ogen).

## Installation

```bash
go get github.com/sst/opencode-sdk-go
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    opencode "github.com/sst/opencode-sdk-go"
)

func main() {
    client, err := opencode.NewClient("http://localhost:3000")
    if err != nil {
        log.Fatal(err)
    }

    // Use the client
    ctx := context.Background()

    // Access the raw ogen client for full API access
    raw := client.Raw()

    // Example API call (adjust based on your actual API)
    // resp, err := raw.ConfigGet(ctx)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // fmt.Printf("Config: %+v\n", resp)
}
```

## Features

- **Type-safe**: Full TypeScript-level type safety in Go using generics
- **Modern**: Built with Go 1.21+ features
- **Fast**: Uses ogen's high-performance JSON codec
- **Flexible**: Raw client access for advanced usage
- **Zero-reflection**: Compile-time code generation

## API Documentation

See the generated client in the `generated/` directory for all available methods and types.

## Development

Short guide to generate, clean, and integrate the Go SDK for local development.

Prepare environment

Install and configure Go and required tools.
Use Go 1.24+ toolchain and ensure GOBIN is set if you install tools into your local bin.

Commands:

- go version
- export GOBIN="$HOME/go/bin" (if needed)
- go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

---

Use go workspaces

Use a go.work file to tie the SDK, TUI and local packages together.
Add local modules with use directives so builds resolve replace paths automatically.

Example workflow:

- From repo root run: go work init ./sdk/go ./packages/tui ./packages/sdk/go
- Run: go work sync

This ensures local replace directives are unnecessary and module resolution is consistent across builds.

---

Clean before fresh build

Always reset generated artifacts and temporary module state before regenerating or rebuilding.
Steps:

- Remove prior generated files: rm -rf sdk/go/types.gen.go sdk/go/generated/\*
- Remove module cache in workspace if needed: go clean -modcache
- Optionally restore SDK folder from .go-cleanup-backup-\* if present
- Re-run go work sync after cleanup

---

Generate OpenAPI spec

Produce a code-first OpenAPI spec from the TypeScript API, then bundle it for generator compatibility.

Commands:

- bun run --conditions=development packages/opencode/src/index.ts generate > schema/openapi.json
- npx redocly bundle schema/openapi.json -o schema/openapi.bundled.json --dereferenced

Notes:

- The bundled, dereferenced JSON is required by oapi-codegen and other generators.
- Verify the spec with npx redocly lint schema/openapi.json

---

Generate Go types

Use oapi-codegen to create types-only Go sources, then keep service layer manual.

Commands:

- go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
- oapi-codegen -generate types -package opencode schema/openapi.bundled.json > sdk/go/types.gen.go

Notes:

- Generate types only to avoid discriminator and complex generation pitfalls.
- Commit types.gen.go only as part of the SDK generation flow, not to hand-edit it.

---

Manual service layer

Copy or implement a lightweight, Stainless-inspired service layer that composes the generated types into usable client services.
Typical structure:

- client.go (Client struct + common request options)
- app.go, session.go, event.go, file.go, find.go (service methods)
- option/ and internal/ helper packages

Reasons:

- Manual services preserve ergonomics, error handling and patterns from Stainless.
- They allow precise control of request options, retries, telemetry, and auth.

---

Differences from Stainless SDK

Key differences to know:

- Generation target: we generate types-only with oapi-codegen, then hand-write services; Stainless historically used a full-code generation approach.
- Discriminator handling: types-only avoids discriminator edge cases that full generators sometimes mishandle.
- Module naming and layout: the generated SDK in this repo uses a local module path (git.j9xym.com/opencode-api-go) and a workspace replace during development.
- Service surface: Stainless SDK might expose different helper primitives and higher-level conveniences not present initially; the hybrid approach gives flexibility but requires implementing some helpers manually.
- Tests and examples: Stainless packages include example clients; mirror those patterns when creating the service layer.

---

Go module and replace strategy

Two approaches for local development:

1. Replace in go.mod (older style)

- packages/tui/go.mod contains:
  replace git.j9xym.com/opencode-api-go => ../../sdk/go

2. Workspace-first (preferred)

- Use go.work to reference all local modules and avoid repeated replace directives.
- go.work is more robust for multi-module workspaces and CI.

Choose workspace-first during active development.

---

Integration with TUI build

How the SDK and TUI cooperate:

- TUI imports the SDK module for API types and client services.
- During local development, the TUI resolves the SDK module via go.work or replace directive to the repo-local sdk/go.
- After generation, ensure sdk/go builds (go build ./...) before attempting to build the TUI to catch SDK regressions early.
- Maintain backward compatibility (or provide a compatibility layer) when SDK surface changes to avoid breaking the TUI.

Practical sequence:

1. Generate bundled OpenAPI spec
2. Generate types: sdk/go/types.gen.go
3. Update or copy manual service files into sdk/go/
4. go work sync (or use replace)
5. cd sdk/go && go mod tidy && go build ./...
6. cd packages/tui && go mod tidy && go build ./...

---

Compatibility and maintenance notes

- API breaks: regenerating the SDK may change types and function signatures. Provide a migration/compatibility layer in sdk/go/compat if breaking changes are expected.
- Tests: run unit tests for both sdk/go and packages/tui after generation.
- Backups: keep .go-cleanup-backup-\* for safe rollbacks when regenerating.
- Linting/formatting: run go fmt and go vet across modules after generation.
- Commit policy: generated artifacts that are intended to be checked in should be clearly documented and regenerated via scripts; avoid manual edits to generated files.

---

Scripts and automation

Add or use repository scripts for reproducible generation:

- scripts/generate-go-sdk.sh — runs spec generation, bundling, oapi-codegen, copies services, updates go.mod, and tidy
- scripts/clean-go-sdk.sh — removes generated files and restores workspace to clean state
- Make CI step run generation check or ensure generated files are up-to-date

Example script outline (high level):

1. generate spec
2. bundle spec
3. oapi-codegen -generate types ...
4. copy services from packages/sdk/go
5. update imports (sed) to new module path
6. cd sdk/go && go mod init <module> && go mod tidy
7. go build ./... && cd ../packages/tui && go build ./...

---

Common troubleshooting

- Missing oapi-codegen: install via go install and ensure GOBIN is in PATH.
- Stale generated types: run clean script and regenerate.
- Import path mismatches: prefer go.work, or run sed to fix import paths inside copied files.
- go mod errors: run go mod tidy in module directory and then go work sync.
- Build failures in TUI post-generation: inspect API surface changes; begin with a compatibility adaptor or update a few core call sites.
