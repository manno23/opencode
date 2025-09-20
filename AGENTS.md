## IMPORTANT

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


## Common Errors / Suggestions:
- Auth failure: Set export GOPRIVATE=github.com/sst/* and add GitHub token to ~/.netrc or use SSH remotes.
- If dev server has runtime issues (e.g., API calls fail), check logs for private repo access.
- If more conflicts arise on merge, use git merge recover-legacy-with-changes -X theirs packages/tui/SDK.md to keep locals.
