# OpenCode Agent Guidelines

## Essential Commands

- **Dev server**: `bun run dev`
- **Type checking**: `bun run typecheck`
- **Generate SDKs**: `bun run generate`
- **Test snapshots**: `bun run test:snapshot`
- **Accept snapshots**: `bun run snapshot:accept`
- **Go SDK test**: `cd sdk/go && go test ./...`
- **Single Go test**: `cd sdk/go && go test -v -run TestName ./pkg`

## Code Style

- **TS/JS**: Prefer functional style, immutable data, `const` over `let`, early returns
- **Go**: `go fmt`, meaningful names, explicit error handling, context for cancellation
- **Formatting**: Prettier (semi: false, printWidth: 120), ES modules with `.js` extensions
- **Avoid**: `else`, `try`/`catch`, `any`, unnecessary destructuring
- **Naming**: camelCase (JS/TS), snake_case (Go vars), PascalCase (Go exports)

## Workflow

1. Understand problem deeply first
2. Investigate codebase thoroughly
3. Plan incrementally, test frequently
4. Run typecheck and snapshot tests before completion

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
