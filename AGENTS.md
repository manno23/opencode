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
- Use as many bun APIs as possible like `Bun.file()`

### TypeScript/JavaScript

- Use TypeScript for all new code
- Prefer functional programming style
- Use immutable data structures where possible
- Follow Prettier formatting (semi: false, printWidth: 120)
- Use ES modules with `.js` extensions in imports
- Prefer `const` over `let`
- Use early returns instead of nested conditionals

### Go

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable names (not single letters unless in loops)
- Handle errors explicitly at each call site
- Use context for cancellation
- Prefer struct composition over inheritance
- Use interfaces for abstraction

### Imports

- Group imports by type (standard library, third-party, local)
- Use absolute imports for internal packages
- Keep import statements organized and sorted

### Error Handling

- Handle errors immediately when they occur
- Return errors from functions rather than panicking
- Use structured logging with context
- Avoid silent error swallowing

### Naming Conventions

- Use camelCase for variables and functions in JS/TS
- Use PascalCase for classes/types in JS/TS
- Use snake_case for variables/constants in Go
- Use PascalCase for exported functions/types in Go
- Prefer descriptive names over abbreviations

## Copilot Instructions

The project includes comprehensive agent instructions in `packages/opencode/src/session/prompt/copilot-gpt-5.txt` that define:

### Workflow

1. Deeply understand the problem before coding
2. Investigate the codebase thoroughly
3. Develop a clear, step-by-step plan
4. Implement changes incrementally
5. Debug systematically
6. Test frequently
7. Iterate until complete

### Communication Style

- Keep responses short and impersonal
- Use warm and friendly yet professional tone
- Sprinkle in light, witty humor when appropriate
- Always communicate clearly and concisely

### Code Search Guidelines

- Prefer workspace symbols over grep for precise identifiers
- Use semantic search for high-level concepts
- Add fully qualified links for referenced symbols
- Include file links for all workspace files

## Development Environment

### Prerequisites

- Bun (package manager and runtime)
- Go 1.24.4+
- Node.js/TypeScript
- Git

### Project Structure

- `packages/opencode/` - Main CLI application
- `packages/sdk/` - SDK generation tools by Stainless
- `packages/tui/` - Terminal user interface
- `sdk/go/` - Generated Go SDK - Will usurp stainless generated sdk
- `sdk/typescript/` - Generated TypeScript SDK - Same as above

### Key Configuration Files

- `tsconfig.json` - TypeScript configuration
- `.editorconfig` - Editor formatting rules
- `bunfig.toml` - Bun configuration
- `go.work` - Go workspace configuration

## Common Development Tasks

### Adding a New Feature

1. Create feature branch from main
2. Implement changes following code style guidelines
3. Add tests for new functionality
4. Update documentation if needed
5. Run `bun run typecheck` and `bun run test:snapshot`
6. Submit pull request

### Debugging Issues

1. Check server logs for API errors
2. Use `slog` for structured logging
3. Test API endpoints with curl/Postman
4. Verify configuration files are valid
5. Check network connectivity to server

### Performance Optimization

1. Profile with built-in tools
2. Minimize API calls and payload sizes
3. Use efficient data structures
4. Implement caching where appropriate
5. Monitor memory usage and garbage collection
