# Incremental Delivery Plan for TUI SDK Migration

## High-Level Strategy

The goal of this plan is to reinstate the legacy Go SDK as the primary client for the TUI, ensuring product velocity and stability, while incrementally introducing an ogen-generated SDK for improved type safety, ergonomics, and maintainability. This approach minimizes risk by avoiding a big-bang migration, allowing for feature flags, rollback paths, and gradual parity validation.

### Key Principles

- **Stability First**: Default to the legacy SDK for all operations in stable releases. Use feature flags to enable ogen implementations per domain group.
- **Facade Pattern**: Introduce stable interfaces in `packages/tui/internal/app/client` that abstract the underlying SDK. TUI UI code depends only on these facades, making switches transparent.
- **Incremental Migration**: Migrate domain groupings (e.g., Sessions, Messages) one at a time, starting with low-risk, high-impact areas. Validate each group with tests and canaries before enabling in production.
- **Spec Workflow**: Maintain a single canonical OpenAPI spec from the server (`openapi.json`), overlay with ogen-specific customizations (`openapi.yaml`), and merge into a unified spec for ogen generation. This keeps the server spec clean while enabling ogen features.
- **Phased Rollout**:
  1. **Phase 1 (Adapters)**: Implement facades and legacy implementations; generate initial ogen client; migrate first group (Sessions).
  2. **Phase 2 (Core Flows)**: Migrate Messages (non-streaming), Tools; add streaming parity.
  3. **Phase 3 (Full Parity)**: Migrate Files, Share, and remaining endpoints; remove flags and legacy code.
  4. **Phase 4 (Optimization)**: Add retries, observability; evaluate committing generated code.
- **Completion Criteria**: All domain groups migrated, flags removed, legacy SDK deprecated, full test coverage, no regressions in latency/error rates, and documentation updated for contributors.

This strategy covers the full path to completion: from dual implementations to a unified ogen-based SDK, with guardrails for risks like spec divergence and streaming differences.

## Domain Groupings for Migration

Prioritize based on change surface and impact:

1. **Health and Metadata**: ping, version, config (low risk, foundational).
2. **Sessions**: CRUD, select/activate (core to TUI state).
3. **Messages**: send, stream, history (high usage, streaming critical).
4. **Tools**: list, invoke, permissions (frequent in agent flows).
5. **Files/Attachments**: upload, list, delete (UI interactions).
6. **Share/Links**: create, fetch (user-facing features).
7. **Admin/Dev Endpoints**: Last, if exposed to TUI.

For each group:

- Add a feature flag (e.g., `TUI_SDK_EXPERIMENT_SESSIONS=true`).
- Defaults off in stable builds; enable progressively in canaries/nightlies.

## Interface and Adapters

Define minimal, stable interfaces per group in `packages/tui/internal/app/client` (e.g., `sessions.go`, `messages.go`):

- **SessionsClient**: `Create(ctx, opts) (*Session, error)`, `List(ctx) ([]*Session, error)`, `Activate(ctx, id string) error`, `Close(ctx, id string) error`.
- **MessagesClient**: `Send(ctx, sessionID, content string, opts) (*Message, error)`, `Stream(ctx, sessionID, content string, opts) <-chan *StreamEvent`.
- **ToolsClient**: `List(ctx) ([]*Tool, error)`, `Invoke(ctx, toolID string, params any) (*ToolResult, error)`.
- **FilesClient**: `Upload(ctx, reader io.Reader, meta FileMeta) (*File, error)`, `List(ctx) ([]*File, error)`, `Delete(ctx, id string) error`.

Implementations:

- **LegacyImpl**: Wraps existing legacy Go SDK (e.g., `next-go` or current client) directly.
- **OgenImpl**: Wraps `packages/tui/api/ogen` client, mapping generated types (e.g., `optional.String` to facade types) and handling differences (e.g., sum types via `IsX/AsX`).

Constructor in `client.go`: Builds per-group impls based on flags, e.g., `newClient(flags) *Client { return &Client{Sessions: if flag { NewOgenSessions() } else { NewLegacySessions() }, ... }`.

TUI code (e.g., `internal/tui/tui.go`) injects the facade client and remains unchanged during migration.

## Spec and Generation Workflow

- **Authoritative Spec**: `packages/tui/api/openapi.json` (generated via `bun run generate` or `opencode generate` from server).
- **Custom Overlay**: `packages/tui/api/openapi.yaml` (manual additions: x-ogen-\*, discriminators for unions, operation grouping).
- **Merged Spec**: `packages/tui/api/unified-openapi.yaml` (build artifact; .gitignore'd) – combine JSON base + YAML deltas using jq/yq or a Go merger.
- **Ogen Config**: `packages/tui/api/ogen.yml` (e.g., `ptr-for-objects: true`, `features: {std-errors: true, client: true}`).
- **Generation**: `packages/tui/api/generate.go` with `//go:generate ogen --clean --config ./ogen.yml --target ./ogen --package api ./unified-openapi.yaml`.
- **Layout**:
  ```
  packages/tui/api/
  ├── openapi.json      # Server-generated (commit if stable)
  ├── openapi.yaml      # Manual deltas (committed)
  ├── unified-openapi.yaml # Merged (gitignore)
  ├── ogen.yml          # Config (committed)
  ├── generate.go       # go:generate (committed)
  └── ogen/             # Generated client (gitignore; generate in CI)
  ```
- **Commands** (via Makefile/package.json scripts):
  - `make spec`: Fetch/generate `openapi.json`.
  - `make merge`: JSON + YAML → `unified-openapi.yaml`.
  - `make sdk`: `go generate ./api/...` (produces `ogen/`).
  - `make build`: `go build -o tui cmd/opencode/main.go`.
  - `make clean`: Remove generated artifacts.

## Streaming and SSE Handling

- **Legacy**: Retain existing streaming path in `LegacyImpl` (e.g., manual HTTP client with event parsing).
- **Ogen**:
  - Use native ogen streaming if supported (e.g., `EventSubscribe(ctx) io.ReadCloser`).
  - Fallback: Thin wrapper in `OgenImpl` using ogen's HTTP client to handle SSE, parsing lines into generated event types (e.g., sum types for Event oneOf).
- Facade unifies: Both impls expose `Stream(...) <-chan *Event` with context cancellation and reconnection logic centralized in `app/client`.

## Auth, Retries, and Observability

- **Centralized in Facade**: `app/client` handles token provider, retries (exponential backoff with jitter), and logging.
- **Ogen-Specific**: Use `WithRequestEditorFn` for auth headers (`Authorization: Bearer <token>`) and telemetry (OpenTelemetry spans for requests).
- **Parity**: Add metrics (latency, errors) to both impls; compare during canaries. Enable tracing for debugging migration issues.

## Risk Management and Guardrails

- **Scoped Changes**: Limit to listed domain groups; defer low-usage endpoints.
- **Flags and Rollback**: Per-group flags; toggle to legacy for instant rollback.
- **SLAs**: No regressions in core flows (session creation <500ms, message streaming <1s latency); monitor error rates.
- **Quality Bars**:
  - Unit tests per impl (mock server).
  - Integration tests for dual paths (call both for smoke validation).
  - Golden tests for event parsing and type mappings.

Key Risks:

- **R-12 (SSE Differences, Prob: 0.3, Impact: High)**: Mitigation – Defer streaming migration; add test harness for events.
- **R-18 (Spec Divergence, Prob: 0.4, Impact: Medium)**: Mitigation – CI diff check between server spec and unified; restrict overlay to x-ogen-\*.
- **R-21 (Dev Friction, Prob: 0.25, Impact: Medium)**: Mitigation – Clear docs/flags; linter to enforce facade usage.

## Decision Records (ADRs)

- **ADR-7**: Legacy Go SDK primary; facade with per-group OgenImpl behind flags (Accepted).
- **ADR-8**: Spec overlay/merge to `unified-openapi.yaml` for ogen (Accepted).
- **ADR-9**: Do not commit generated ogen client; generate in CI (Proposed).

## Documentation and Contributor Guidance

- **Update `SDK.md`**:
  - Describe facade pattern with code snippets (interfaces, constructors).
  - Document flags (e.g., env vars) and enabling OgenImpl.
  - Detail spec merge (prereqs: jq/yq; commands via Makefile).
  - Migration checklist per group: type mappings, optionals → `optional.*`, unions → sum types, errors → toasts.
- **High-ROI Refactors**: Centralize API calls in `app/app.go`; remove manual JSON in `tui.go`; add error adapter.
- **Useful Files**: `internal/tui/tui.go` (UI orchestration), `internal/app/app.go` (API usage).

## CI/CD Integration

- **Pipeline Stages**:
  - Generate spec → Merge → SDK gen → Build TUI → Lint/Typecheck → Unit/Integration tests.
  - Canary Branches: Dual-path smoke tests (e.g., Sessions.Create via both impls).
- **Gates**: PRs link to milestones; pass tests before enabling flags in mainline.

## Delivery Milestones

- **P1-Adapters (ETA: 2025-09-24)**: Facades for Sessions/Messages; LegacyImpl.
- **Ogen Client Gen (ETA: 2025-09-26)**: Unified spec; initial generation.
- **Sessions via Ogen (ETA: 2025-09-30)**: First group behind flag; smoke tests.
- **Messages Non-Stream (ETA: 2025-10-03)**: Send/history migration.
- **Canary on Core (ETA: 2025-10-07)**: Enable flags in nightly; monitor.

## Next 7 Days Actions

1. Add facade interfaces/constructors in `internal/app/client`.
2. Implement LegacyImpl for Sessions/Messages.
3. Fix `api/generate.go` to use `unified-openapi.yaml`; validate ogen.yml.
4. Wire Makefile/scripts (spec, merge, sdk, build); update `SDK.md`.
5. Generate ogen client; build TUI; smoke-test Sessions via OgenImpl.
6. Document flags/rollback in `SDK.md`; add .gitignore for artifacts.

## Quality Gates

- PRs reference milestones and assign risk owners.
- Escalate if burn-down >15% (e.g., defer streaming).
- Review `SDK.md` with every code change for accuracy.
- Full parity: 100% test coverage, no flags, legacy code removed.

=============================================================================================

# Thoughts on what to keep, and what can be discarded.

Here’s a pragmatic map for your SDK refactor: what to lean on, what to migrate into the generated client, and what you can likely delete or simplify once ogen
covers it.

## Keep in the app (UI/UX-specific)

These are TUI concerns and should remain outside the SDK:

• Bubble Tea components and layout: editor/messages/modals/toasts, debounced key handling, leader sequences
• Files: packages/tui/internal/tui/tui.go, internal/components/\*, internal/layout, internal/styles, internal/theme, internal/components/{chat,dialog,status,toast}
• Command routing and keybind logic
• Files: internal/commands, usage in tui.go executeCommand branch
• App state and UI state orchestration (which modal to show, where to place overlays, cursor math, ANSI conversion)
• Files: internal/app for UI state, internal/util for TUI helpers (place overlay, ANSI conversion)

## Move into the SDK (generator-friendly, API-shaped)

These should be centralized in the generated client or a thin hand-written layer next to it:

• API request/response typing and union handling
• Replace all opencode.F optional wrappers with ogen-generated Opt* types and helpers
• Use SDK union discriminators directly; avoid local re-wrapping of unions
• Endpoint coverage
• Port all calls listed in packages/sdk/next-go/SDK.md (Session*, Project*, File*, Find*, CommandList, ConfigProviders, Tool* and TUI API endpoints) into the ogen
client
• Ensure streaming (EventSubscribe) returns a Reader-like handle with clean close/cancel via context
• Pagination, filtering, and params shaping
• Any repeated construction of params in internal/app or components can be collapsed into SDK request builders
• Permissions API
• Provide a typed Permissions.Respond method (your SDK.md names it PostSessionIdPermissionsPermissionID) with clean enums
• Error normalization
• Map API error payloads into typed errors (ProviderAuthError, UnknownError, etc.) so TUI can type-switch without inspecting raw payloads

Concretely, the following code paths can become thin SDK calls (remove custom glue):

• Session interactions in TUI and components:
• Session.Get, Children, Share/Unshare, Revert/Unrevert, Summarize, Init, Prompt, Command, Shell, Update, Delete, Abort, List, Messages, Message
• References:
• packages/tui/internal/tui/tui.go: calls around lines ~120–450, ~650–750, ~1215–1475
• packages/opencode/internal/app/app.go: all Session.\* usage noted in SDK.md
• packages/opencode/internal/components/chat/messages.go: Message, Revert, Unrevert paths

• App endpoints:
• App.Providers, App.Init, App.Log
• Move struct shapes and request composition into SDK
• TUI “API” commands (local HTTP-ish control messages in tui.go):
• /tui/open-help, /tui/open-sessions, /tui/open-timeline, /tui/open-themes, /tui/append-prompt, /tui/submit-prompt, /tui/clear-prompt, /tui/execute-command,
/tui/show-toast
• For “append/submit/clear/execute/show-toast,” you already spec’d corresponding SDK endpoints (TuiAppendPrompt, TuiSubmitPrompt, TuiClearPrompt,
TuiExecuteCommand, TuiShowToast). Route these through the SDK and drop the custom JSON decoding in tui.go. The UI can just call the SDK methods.

## Simplify or delete after moving to SDK

• Optional field helpers
• Remove opencode.F in app/components once ogen Opt* builders exist
• Manual JSON marshalling for API.Request in tui.go
• Replace the switch handling of msg.Path bodies with typed SDK calls; delete the body structs in tui.go and let the generator’s request types own it
• Local “response union” switches where the SDK can return typed results
• E.g., ToolIdsRes, ToolListRes: if ogen emits tagged unions and helpers, switch on SDK type instead of local intermediates
• Event stream plumbing that’s duplicative
• If your SDK exposes EventSubscribe returning a stream + helper to parse events into strongly typed opencode.EventListResponse* values, you can:
• Keep the high-level event handling in the app/TUI
• Drop any local SSE wire format handling, retry logic, or demux that the SDK can own

## Minimal app-side that you still need

• Mapping SDK model objects into UI Messages
• app.Message wraps SDK unions for rendering; keep the minimal wrapper, but remove any data massaging that the SDK can do (timestamps, IDs, share URL locations)
• Lightweight adapters for UX-only flows
• The logic in tui.go for switching from child session to parent before sending commands is UX state; keep that outside the SDK, but keep the calls themselves
thin

## Pitfalls to watch out for when swapping SDKs

• Enum names and param names
• Example: SessionPermissionRespondParamsResponse\* enums might be renamed or nested differently in the ogen output; centralize these in a small adapter layer to
avoid sprinkling renames across the TUI
• Pointer vs value returns
• Your SDK.md shows mix of pointers and slices; update nil checks accordingly
• Union types
• Ensure ogen generates usable Go unions (sum types) and expose ergonomic helpers (AsX, IsX). If it doesn’t, add tiny wrappers next to the SDK, not in the TUI
• Streaming cancellation
• Confirm Context cancellation on EventSubscribe fully tears down the Reader; if not, provide an SDK-level Close

## High-ROI targets

• packages/opencode/internal/app/app.go
• Centralizes most API interactions; migrate all request construction to SDK; reduce to orchestrating calls and UI state
• packages/tui/internal/components/chat/messages.go
• Message loading, revert/unrevert details, tool/text/file part handling; replace API shaping with SDK calls and use typed unions
• packages/tui/internal/tui/tui.go
• Remove the API.Request body decoding and method branching for /tui/\* paths; call the SDK directly
• packages/opencode/internal/util/apilogger.go
• Replace with SDK App.Log

After this shift:

• The SDK owns all HTTP shape: params, unions, errors, streaming
• The TUI only orchestrates UI and calls typed SDK methods
• You shed custom option wrappers, manual JSON decoding for “/tui/\*”, and any local type aliases that mirror API responses

If you want, I can produce a concrete diff plan for app.go, messages.go, and tui.go to swap in the ogen client, listing exact call-site substitutions and
enum/param renames.

==============================================================================================
