# OpenCode Agent Guidelines

## Build/Lint/Test Commands

### Root Package Commands

- **Development server**: `bun run dev`
- **Type checking**: `bun run typecheck`
- **Generate SDKs**: `bun run generate`
- **Snapshot testing**: `bun run test:snapshot`
- **Accept snapshots**: `bun run snapshot:accept`

### SDK Generation Commands

- **Generate all SDKs**: `bun run generate:all`
- **Generate TypeScript SDK**: `bun run generate:ts`
- **Generate Go SDK**: `bun run generate:go`
- **Clean generated files**: `bun run clean`
- **Test SDK generation**: `bun run test`
- **Test TypeScript SDK**: `bun run test:ts`
- **Test Go SDK**: `bun run test:go`

### Individual Package Commands

- **TypeScript packages**: `bun run typecheck` (in each package)
- **Go SDK testing**: `cd sdk/go && go test ./...`
- **Single test file**: `cd sdk/go && go test -v ./path/to/package`
- **Single test function**: `cd sdk/go && go test -v -run TestFunctionName ./path/to/package`

## Code Style Guidelines

### General Principles

- Try to keep things in one function unless composable or reusable
- DO NOT do unnecessary destructuring of variables
- DO NOT use `else` statements unless necessary
- DO NOT use `try`/`catch` if it can be avoided
- AVOID `try`/`catch` where possible
- AVOID `else` statements
- AVOID using `any` type
- AVOID `let` statements
- PREFER single word variable names where possible
- Use as many bun apis as possible like Bun.file()

## Debugging

- To test opencode in the `packages/opencode` directory you can run `bun dev`
