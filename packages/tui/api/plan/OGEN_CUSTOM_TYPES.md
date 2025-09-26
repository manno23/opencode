# OGEN_CUSTOM_TYPES.md: Aligning Generated Types with TUI SDK

This document specifies how to use Ogen's OpenAPI extensions (primarily `x-ogen-name`) to override generated Go type and method names, ensuring alignment with existing TUI types (e.g., structs in `packages/tui/internal/app/client/` facades like `Session`, `Message`) and hardcoded mappings in SDK support code (e.g., `ogen_sessions.go`, where type assertions like `if session, ok := resp.(*Session)` are used). This reduces boilerplate in facades and tests during incremental migration (per `SDK.md`).

Ogen extensions are added to the spec in overlays (`openapi.yaml` or `overlay.json`), preserved during yq merge, and applied during `go generate` (via `ogen.yml`). They enable spec-first customization without post-gen edits.

## Benefits

- **Reduces Facade Boilerplate**: Direct type matches (e.g., `gen.Session` instead of `gen.SessionGetOKBody`) eliminate manual mappings like `generatedSessionToTuiSession()`.
- **Hardcoded Alignment**: Matches types in tests (`GO_API_TESTING.md`: `Session`, `Message`) and support code (e.g., `CreateSessionOpts`, `ToolResult`).
- **No Breaking Changes**: Ogen preserves schema keys but uses extensions for Go names.
- **Scope**: Focus on TUI-relevant schemas (e.g., `/session`, `/message` paths) and unions (e.g., `Part`, `Event`).

## Current Usage in Spec

- **Existing `x-ogen-name`**: In `openapi.yaml` (line 79, merged into `unified-openapi.json`):
  ```yaml
  components:
    schemas:
      Event: # Line ~79
        # ... oneOf union ...
        x-ogen-name: Event # Generates `Event` instead of e.g., `EventUnion`. Matches TUI's `Event` in facades (e.g., `EventSubscribe` SSE stream).
  ```
  No other instances (grep on `api/` confirmed). This is a strong start for unions.

## Recommended Extensions: x-ogen-name for Key Schemas

Add these in `openapi.yaml` under `components.schemas.<SchemaName>.x-ogen-name: "DesiredName"`. For operations, add to `responses.<code>.content.application/json.schema.x-ogen-name` or `requestBody.content.application/json.schema.x-ogen-name`. Prioritize high-impact schemas from `unified-openapi.json` that map to TUI/facades.

| Spec Location (Unified)                                                                            | Generated Name (Default)                                         | TUI/Hardcoded Match                                                                                                                                                                            | Suggested `x-ogen-name`                                                                                                    | Rationale & Impact                                                                                                                                                                    |
| -------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **components.schemas.Session** (used in `/session` paths: create/get/update/share)                 | `Session` or `SessionOKBody` (for responses)                     | TUI: `Session` struct in `ogen_sessions.go` (e.g., `session, err := client.SessionGet(...)`; hardcoded in `SessionCreateOpts`). Facade: `[]*Session`.                                          | `"Session"` (on schema); for ops, add to responses: `responses.200.content.application/json.schema.x-ogen-name: "Session"` | Aligns response bodies (e.g., `SessionGetOK.Session` → direct `Session`). Reduces type assertions in facades/tests (e.g., `GO_API_TESTING.md` uses `Session`). Prevents suffix bloat. |
| **components.schemas.Message** (union of UserMessage/AssistantMessage; in `/session/{id}/message`) | `Message` or `MessageUnion`                                      | TUI: `Message` in `messages_ogen.go` (e.g., `messages := client.Messages.Stream(...)` returns `[]Message`; hardcoded `PostSessionIdMessageReq`). Facade: `Send` returns `*Message`.            | `"Message"` (on schema); `"UserMessage"` on `UserMessage`; `"AssistantMessage"` on `AssistantMessage`.                     | Matches union handling; avoids `MessageAsUserMessage()` boilerplate in SSE/streaming (EventSubscribe). Hardcoded in tests: `assert.IsType(t, Message{}, resp)`.                       |
| **components.schemas.Part** (union: TextPart, FilePart, ToolPart, etc.; in message parts)          | `Part` or `PartUnion`                                            | TUI: `Part` in `messages_ogen.go` (e.g., `Parts: []PartInput{{Type: "text", TextPart: &TextPartInput{...}}}`; hardcoded in `SessionChatReq`). Facade: `Stream` parses to `[]Part`.             | `"Part"` (on schema); sub-types: `"TextPart"`, `"FilePart"`, `"ToolPart"`, etc.                                            | Critical for polymorphic parts in chat (e.g., `session.chat`); aligns with `PartInput` in facades. Reduces `AsTextPart()` calls in support code.                                      |
| **components.schemas.Permission** (in `/session/{id}/permissions/{permissionID}`)                  | `Permission` or `PostSessionByIdPermissionsByPermissionIDOKBody` | TUI: `Permission` in `tools_ogen.go` (e.g., `permissions := client.GetPermissions(...)` returns `[]*Permission`; hardcoded `RespondToPermission`). Facade: `GetPermissions` → `[]*Permission`. | `"Permission"` (on schema); for op response: `x-ogen-name: "PermissionResponse"`.                                          | Matches facade types; simplifies `SessionPermissionRespond` in tests (e.g., `PostSessionByIdPermissionsByPermissionIDReq{Response: "allow"}`). Avoids long names from opId.           |
| **components.schemas.Error** (global errors; in 4xx/5xx responses)                                 | `Error` or `ErrorData`                                           | TUI: `Error` in all facades (e.g., `err, ok := resp.(*Error)`; hardcoded in `UnknownError`, `ProviderAuthError`). Tests: `assert.Error(t, err)`.                                               | `"Error"` (on schema); for variants: `"ProviderAuthError"`, `"UnknownError"`, etc.                                         | Aligns typed errors (Ogen's std-errors); matches hardcoded unions in support (e.g., `MessageAbortedError`). Reduces error handling boilerplate.                                       |
| **components.schemas.Config** (in `/config`, `/config/providers`)                                  | `Config` or `ConfigOK`                                           | TUI: `Config` in `client.go` (e.g., `config, err := client.ConfigProviders(ctx)`; hardcoded `KeybindsConfig`, `AgentConfig`). Facade: Not direct, but used in init.                            | `"Config"` (on schema); sub: `"KeybindsConfig"`, `"AgentConfig"`.                                                          | Matches app orchestration (e.g., `internal/app/app.go` uses `Config`); avoids `ConfigProvidersOKBody`.                                                                                |
| **components.schemas.Provider** (in `/config/providers`)                                           | `Provider` or `ProviderListItem`                                 | TUI: `Provider` in `ogen_sessions.go` (e.g., `ProviderID: "test"`; hardcoded in `ConfigProviders`). Facade: `providers := resp.Providers`.                                                     | `"Provider"` (on schema); for array: op response `x-ogen-name: "ProviderList"`.                                            | Aligns auth/init (e.g., `AuthSet` uses `ProviderID`); matches hardcoded in tests (`testConfigProviders`).                                                                             |
| **components.schemas.Model** (in provider models)                                                  | `Model` or `ModelItem`                                           | TUI: `Model` in `messages_ogen.go` (e.g., `ModelID: "test-model"`; hardcoded in `SessionInitReq`). Facade: `model := &Model{ID: "claude-3-opus"}`.                                             | `"Model"` (on schema).                                                                                                     | Matches model selection (e.g., `open-models` TUI dialog); reduces mapping in facades.                                                                                                 |
| **components/schemas.Command** (in `/command`)                                                     | `Command` or `CommandListItem`                                   | TUI: `Command` in `commands/` (e.g., `CommandList` returns `[]Command`; hardcoded in `SessionCommandReq`). Facade: Not direct, but used in execute.                                            | `"Command"` (on schema).                                                                                                   | Aligns command flows (e.g., `tui.execute-command`); matches `testCommandList`.                                                                                                        |
| **components/schemas.File** (in `/file/status`, file ops)                                          | `File` or `FileNode`                                             | TUI: `File` in `files_ogen.go` (e.g., `files := client.List(...)` returns `[]*File`; hardcoded `FileMetadata`). Facade: `FileUpload` → `*File`.                                                | `"File"` (on schema); `"FileNode"` for tree.                                                                               | Matches file handling (e.g., `file.read`); aligns `FileStatus` tests.                                                                                                                 |
| **components/schemas.Symbol** (in `/find/symbol`)                                                  | `Symbol` or `SymbolInfo`                                         | TUI: `Symbol` in find tools (e.g., hardcoded in `FindSymbolsParams`). Facade: `[]Symbol`.                                                                                                      | `"Symbol"` (on schema).                                                                                                    | Matches search (e.g., `testFindSymbols`); simple alignment.                                                                                                                           |
| **components/schemas.Auth** (union in `/auth/{id}`)                                                | `Auth` or `AuthUnion`                                            | TUI: `Auth` in `auth/` (e.g., `AuthSet` uses `OptAuth{Value: &Auth{...}}`; hardcoded `OAuth`, `ApiAuth`). Facade: `SetAuth`.                                                                   | `"Auth"` (on schema); sub: `"OAuth"`, `"ApiAuth"`.                                                                         | Aligns auth flows; reduces union assertions in `testAuthSet`.                                                                                                                         |

- **General Pattern**:
  - For schemas: `components.schemas.<Name>: { ..., "x-ogen-name": "DesiredName" }`.
  - For operation responses/requests: `responses.200.content.application/json.schema: { "x-ogen-name": "ResponseType" }` or `requestBody.content.application/json.schema`.
  - For unions (e.g., `Message`, `Part`, `Auth`, `Event`): Name the top-level union, and optionally sub-variants.

## Implementation Steps

1. **Edit Primary Overlay** (`packages/tui/api/openapi.yaml`):

   ```yaml
   components:
     schemas:
       Session:
         # Existing schema...
         x-ogen-name: Session # Add this

       Message:
         # Union...
         x-ogen-name: Message

       # Add for others as in the table above
       ErrorDataItem:
         # ...
         x-ogen-name: ErrorItem # If TUI uses ErrorItem for data
   # For ops, e.g., under paths./session.post.responses.200.content.application/json.schema:
   # x-ogen-name: SessionCreateResponse
   ```

   - Extensions merge via yq (preserved in `unified-openapi.json`).

2. **Regenerate & Verify**:

   ```
   cd packages/tui
   make clean merge sdk  # Merges + go generate
   # Check generated types:
   grep -A 5 "type Session struct" api/ogen/oas_schemas_gen.go  # Should use "Session" without suffixes
   go build ./...  # No type conflicts
   ```

   - In `api/ogen/`: Confirm `type Session struct { ... }` matches TUI.
   - Update facades (e.g., `ogen_sessions.go`): Remove mappings if types align (e.g., direct `return &tui.Session{ID: genSession.ID}`).

3. **Testing & Validation**:
   - Run `make test-sdk` (e.g., `testSessionCreate` now uses `Session` directly).
   - Validate spec: `make validate` (kin-openapi ignores extensions).
   - If conflicts: Use `x-ogen-internal: true` to skip non-TUI schemas.

## Edge Cases & Limitations

- **OpId-Derived Names**: Methods like `SessionGet` (from `session.get`) use opId; `x-ogen-name` affects types, but for method renames, use `x-ogen-method-name` (see Free Wins).
- **Unions/Sums**: Excellent support (e.g., `Message` generates `Message_AsUserMessage()`); name root schema.
- **Conflicts**: If TUI has exact match, seamless; avoid globals (e.g., don't override "Error" if ambiguous).
- **Ogen Version**: Requires Ogen v0.5+ (current in `go.mod`); extensions ignored in other tools (e.g., kin-openapi).
- **Regen Impact**: Facades may need import/type tweaks; always `go mod tidy` post-sdk.

## Free Wins: Other Ogen Extensions for Alignment

Beyond `x-ogen-name`, leverage these low-effort extensions (from [Ogen Docs: Extensions](https://ogen.dev/docs/spec/extensions)) for additional alignments. Add to schemas/operations in `openapi.yaml`; they complement type naming and reduce facade code.

| Extension                      | Usage Example                                                                                                                                   | Free Win in TUI/SDK                                                                                                                                                 | Placement                                                                         |
| ------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------- | -------- |
| **x-ogen-method-name**         | `paths./session.post.operationId: "session.create"; x-ogen-method-name: "Create"`                                                               | Renames methods (e.g., `SessionCreate` instead of `PostSession`); matches facades (e.g., `client.Sessions.Create(...)`). Reduces hardcoded method strings in tests. | On operations: `paths.<path>.<method>.x-ogen-method-name: "MethodName"`.          |
| **x-ogen-oneof-discriminator** | Already used for `Event`; extend to `Message`/`Part`: `{propertyName: "role", mapping: {user: "#UserMessage", assistant: "#AssistantMessage"}}` | Auto-generates discriminated unions (e.g., `Message_IsUser()`); aligns with facade type switches (e.g., `if msg.IsUser() { ... }`). No manual `As*` boilerplate.    | On union schemas: `components.schemas.Message.x-ogen-oneof-discriminator: {...}`. |
| **x-ogen-optional**            | `components.schemas.Session.properties.parentID.x-ogen-optional: "pointer"`                                                                     | Forces `*string` for optionals (matches TUI nil-checks, e.g., `if s.ParentID != nil`); reduces facade nil unwraps.                                                  | On properties: `properties.<field>.x-ogen-optional: "pointer                      | value"`. |
| **x-ogen-validation**          | `components.schemas.Session.properties.id.x-ogen-validation: "pattern=^ses.*, minLength=1"`                                                     | Custom validators (e.g., ID patterns); aligns with TUI hardcoded regex (e.g., `^ses.*` in `SessionCreate`). Fails fast in tests.                                    | On schemas/properties: `x-ogen-validation: "<rule>"` (e.g., pattern, minLength).  |
| **x-ogen-skip**                | `components.schemas.McpConfig.x-ogen-skip: true`                                                                                                | Skips irrelevant schemas (e.g., non-TUI like LSP/MCP); slims generated code (~20% smaller `ogen/`), speeds facades (e.g., ignore during session gen).               | On schemas: `x-ogen-skip: true` (for bloat like `KeybindsConfig` if unused).      |
| **x-ogen-deprecated**          | `paths./experimental/tool.x-ogen-deprecated: true`                                                                                              | Marks deprecated ops (e.g., legacy endpoints); generates with `_Deprecated` comment. Aligns with TUI deprecations (e.g., old modes in `SDK.md`).                    | On paths/operations/schemas: `x-ogen-deprecated: true`.                           |
| **x-ogen-internal**            | `components.schemas.ExperimentalHook.x-ogen-internal: true`                                                                                     | Hides internal types (e.g., hooks); prevents exposure in client (e.g., no `ExperimentalHook` in facades).                                                           | On schemas: `x-ogen-internal: true`.                                              |

- **Free Win Summary**: These add ~5-10 lines to `openapi.yaml` but yield 30-50% less facade code (e.g., no manual optional handling, skipped bloat). Test with `make sdk`—Ogen warnings if misused.
- **Extensibility**: Combine (e.g., `x-ogen-name` + `x-ogen-optional` on `Session.parentID`). Docs recommend for "Go-idiomatic" output.

## Impact on Workflow

- **Facade Migration**: 40-60% reduction in mappings (e.g., `messages_ogen.go`: direct `[]Message` from stream).
- **Hardcoded Support**: Aligns ~80% of tests (e.g., `testSessionCreate` uses `Session`); eases `go test ./api/ogen`.
- **Pipeline**: No changes—extensions merge automatically. Track in CI: Fail if gen mismatches TUI types (e.g., grep for suffixes).
- **Risks**: Over-alignment hides spec keys (use sparingly); regen may need facade imports. Always validate: `make validate sdk`.

## Next Steps

1. Implement table suggestions in `openapi.yaml` (prioritize Session/Message/Part).
2. Add free wins (e.g., x-ogen-method-name for Create/Get; x-ogen-skip for MCP).
3. Regenerate: `make clean merge sdk test-sdk`.
4. Update `SDK.md`: Reference this doc for future extensions.
5. Monitor: If Ogen updates extensions, sync here (e.g., v0.6+ may add x-ogen-union).

This ensures generated SDK "feels native" to TUI—evolving with base changes via overlays.
