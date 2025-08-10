## Build/Lint/Test Commands

- `bun run typecheck` - Type check all packages
- `bun test` - Run all tests
- `bun test <pattern>` - Run specific test files
- `bun test --test-name-pattern <regex>` - Run specific tests by name
- `bun run dev` - Start development mode
- `bun run build` - Build all packages

## Code Style Guidelines

- **Formatting**: 2 spaces, 120 chars, no semicolons (Prettier config)
- **Types**: No `any`, prefer strict types, use `zod` for validation
- **Naming**: Single-word variables, camelCase for functions/vars, PascalCase for types
- **Imports**: Use `import type` for types, avoid default exports
- **Error handling**: Avoid try/catch unless necessary, prefer early returns
- **Functions**: Keep single-purpose, avoid else statements, prefer composition
- **Variables**: Use `const` only, no `let`
- **APIs**: Prefer Bun APIs (Bun.file, Bun.test, etc.) over Node.js equivalents
