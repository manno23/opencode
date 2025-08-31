## Plan
  
1. Understand the current TUI build flow
    Review packages/opencode/script/publish.ts to see how each platform/arch bundle is created: the script cross‑compiles the Go TUI (go build) and then bundles the TypeScript CLI with Bun for each target before packaging and publishing.

2. Catalogue existing API usage and generated code
    Inspect the OpenAPI spec (openapitools.json) and the generated 3rd‑party client code to list every endpoint, model, and helper currently relied on by the TUI.
    Identify features we must keep (auth, retries, error types, pagination, etc.).

3. Design an internal client library
    Decide on language (TypeScript) and runtime assumptions (Bun/Node).
    Define minimal HTTP wrapper (using fetch) and type definitions for requests/responses.
    Plan for code organization (e.g., packages/opencode/src/api), error handling, and serialization.

4. Implement a custom generator or hand‑written client
    Option A: Write a small script that parses the OpenAPI spec and emits TypeScript types and endpoint functions.
    Option B: Hand‑craft the client using the spec as reference.
    Ensure the output is deterministic so it can be committed and reviewed like normal code.

5. Integrate with the TUI build/publish pipeline
    Import the new client library into the TUI source.
    Confirm bun build (as used in publish.ts) bundles the client correctly for all targets.
    Remove references to the third‑party generated client in build scripts and dependencies.

6. Testing and validation
    Add unit tests for the client functions.
    Run bun test at repo root and smoke‑test the built binaries by executing ./dist/<pkg>/bin/opencode --version during publish (already done in publish.ts).
    Verify type safety with pnpm -w run typecheck.

7. Gradual rollout
    Behind a feature flag, switch the TUI to use the internal client while keeping the generated one as fallback.
    Collect feedback, fix issues, and eventually remove the old generated client.

8. Cleanup and documentation
    Delete openapitools.json and any generator scripts or configs.
    Update README.md/developer docs to explain how the new client is maintained and how to regenerate types if the API spec changes.



## Refinement of Step 4 – Implement a custom generator or hand‑written client

1. Leverage the server’s own OpenAPI emitter
    Call the existing openapi() helper in packages/opencode/src/server/server.ts to generate a JSON spec from the Zod route definitions.
    This keeps the client code aligned with whatever structure the server currently exposes without having to maintain a separate spec file.

2. Handle streaming and other non‑standard responses manually
    The spec produced above will describe basic request/response shapes but won’t fully capture the streamSSE endpoints used for event subscriptions.
    For such cases, add bespoke client helpers (e.g., subscribeEvents()) that use EventSource/fetch with streaming response bodies instead of relying purely on the spec.

3. Generate types and thin request wrappers
    Write a small script that reads the spec and emits:
    TypeScript interfaces for request/response bodies and path parameters.
    A lightweight HTTP wrapper that performs fetch calls and parses JSON.
    Keep the generator deterministic so the emitted code can be committed to the repo.

4. Layer manual extensions where the spec falls short
    Because the API has grown organically, expect idiosyncrasies—optional fields, ad‑hoc error envelopes, or endpoints whose behavior isn’t fully described in OpenAPI. After code generation, review each endpoint and augment the generated client with:
    Custom error types or unions when the spec uses any/object.
    Convenience methods that mirror current stainless features (retries, pagination helpers, auth token injection, etc.).

5. Validate parity against existing stainless SDK
    Diff the newly generated client with the current stainless output to surface missing models or endpoints.
    Add smoke tests that invoke each method against a mocked server to ensure request shapes and parsing logic remain compatible with the server’s organic evolution.

6. Iterate alongside server changes
    Update the generator script whenever new Zod schemas or routes are added. Because the generator derives directly from server.ts, the client remains in lockstep with the server’s organically evolving API, eliminating the “re‑run stainless” step noted in the package guidelines.



packages/opencode/packages.json


src/plugin/index.ts:6:10 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'createOpencodeClient'.

6 import { createOpencodeClient } from "@opencode-ai/sdk"
           ~~~~~~~~~~~~~~~~~~~~



src/plugin/index.ts:17:21 - error TS7019: 
         Rest parameter 'args' implicitly has an 'any[]' type.

17       fetch: 
         async (...args) => Server.app().fetch(...args),        ~~~~~~~




src/plugin/index.ts:17:52 - error TS2556: 
         A spread argument must either have a tuple type or be passed to a rest parameter.

17       fetch: 
         async (...args) => Server.app().fetch(...args),                                       ~~~~~~~




src/server/routes.ts:2:27 - error TS2307: 
         Cannot find module '../util/validation' or its corresponding type declarations.

2 import { validator } from "../util/validation";             ~~~~~~~~~~~~~~~~~~~~





src/server/routes.ts:4:14 - error TS2451: 
         Cannot redeclare block-scoped variable 'withValidation'.

4 export const withValidation = (schema: 
         z.ZodObject, handler: any) => async (ctx) => {~~~~~~~~~~~~~~


src/server/routes.ts:4:40 - error TS2707: 
         Generic type 'ZodObject<T, UnknownKeys, Catchall, Output, Input>' requires between 1 and 5 typ
e arguments.

4 export const withValidation = (schema: 
         z.ZodObject, handler: any) => async (ctx) => {                          ~~~~~~~~~~~



src/server/routes.ts:4:77 - error TS7006: 
         Parameter 'ctx' implicitly has an 'any' type.

4 export const withValidation = (schema: 
         z.ZodObject, handler: any) => async (ctx) => {                                                               ~~~



src/server/routes.ts:9:20 - error TS18046: 
         'error' is of type 'unknown'.

9     ctx.throw(400, error.errors || "Validation failed");      ~~~~~



src/server/routes.ts:12:14 - error TS2451: 
         Cannot redeclare block-scoped variable 'withValidation'.

12 export const withValidation = (schema: 
         z.ZodObject, handler: any) => async (ctx) => { try { const parsed = validator(schema, ctx); re
turn handler({ ... ~~~~~~~~~~~~~~



src/server/routes.ts:12:40 - error TS2707: 
         Generic type 'ZodObject<T, UnknownKeys, Catchall, Output, Input>' requires between 1 and 5 ty
pe arguments.

12 export const withValidation = (schema: 
         z.ZodObject, handler: any) => async (ctx) => { try { const parsed = validator(schema, ctx); re
turn handler({ ...                                      ~~~~~~~~~~~



src/server/routes.ts:12:77 - error TS7006: 
         Parameter 'ctx' implicitly has an 'any' type.

12 export const withValidation = (schema: 
         z.ZodObject, handler: any) => async (ctx) => { try { const parsed = validator(schema, ctx); re
turn handler({ ...                                                                           ~~~



src/server/routes.ts:13:50 - error TS18046: 
         'error' is of type 'unknown'.

13 ctx, parsed }); } catch (error) { ctx.throw(400, error.errors || "Validation failed"); } };                                                ~~~~~

../plugin/src/index.ts:2:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'Event'.

2   Event,~~~~~

../plugin/src/index.ts:3:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'createOpencodeClient'.

3   createOpencodeClient,~~~~~~~~~~~~~~~~~~~~

../plugin/src/index.ts:4:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'App'.

4   App,~~~

../plugin/src/index.ts:5:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'Model'.

5   Model,~~~~~

../plugin/src/index.ts:6:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'Provider'.

6   Provider,~~~~~~~~

../plugin/src/index.ts:7:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'Permission'.

7   Permission,~~~~~~~~~~

../plugin/src/index.ts:8:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'UserMessage'.

8   UserMessage,~~~~~~~~~~~

../plugin/src/index.ts:9:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'Part'.

9   Part,~~~~

../plugin/src/index.ts:10:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'Auth'.

10   Auth, ~~~~

../plugin/src/index.ts:11:3 - error TS2305: 
         Module '"@opencode-ai/sdk"' has no exported member 'Config'.

11   Config, ~~~~~~


Found 22 errors in 3 files.

Errors  Files 3  src/plugin/index.ts:6 9  src/server/routes.ts:210  ../plugin/src/index.ts:2

