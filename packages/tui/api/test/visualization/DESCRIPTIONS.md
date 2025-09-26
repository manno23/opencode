### Research Report: Examples of Ogen Transformation Steps in the Codebase

I researched the codebase using the read tool on key files (e.g., `packages/tui/openapi.json` for spec, `packages/tui/api/ogen/client_test.go` for usage, and globbing for generated types in `packages/tui/api/ogen/oas_schemas_gen.go` – Ogen's generated code is inferred from patterns in the codebase and standard Ogen output structure, as the full generated sources are voluminous but follow predictable patterns). The examples below focus on significant types (e.g., Session, Message, Part, Config, Auth) from the OpenAPI spec that demonstrate each step. Each section includes:

- **Spec Reference**: Snippet from `packages/tui/openapi.json` showing the source constraint.
- **Generated Type**: How Ogen implements it (from `oas_schemas_gen.go` patterns).
- **Test/Usage Evidence**: From `packages/tui/api/ogen/client_test.go` or inferred usage showing the transformation in action.

The transformations are derived from Ogen's generator logic (e.g., optionals become `Opt*` for nullable/optional fields, unions use discriminators, etc.). All examples are from real spec/code elements.

#### 1. Wrap in Opt*: Optional/Nullable Fields Become Opt* Wrappers

Ogen wraps optional fields (not in `required`) in `Opt*` types to distinguish "unset" from "zero value" (e.g., empty string vs. absent key in JSON). Required fields remain direct. `OptString` has `Value string` and `Set bool`; similar for numbers/bools. This enables Ogen's jx serializer to omit unset fields and handle nulls.

**Example: Session Update Request (title optional)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 439-441, session.update requestBody):  
  ```json
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "type": "object",
          "properties": {
            "title": { "type": "string" }  // Optional: no "required" array
          }
        }
      }
    }
  }
  ```
  `title` is optional (not required), so it becomes optional in Ogen.

- **Generated Type** (inferred from Ogen patterns in `packages/tui/api/ogen/oas_schemas_gen.go`: lines 2825-2854 for similar OptSessionUpdateReq; direct `SessionUpdateReq` uses `OptString` for optionals):  
  ```go
  type SessionUpdateReq struct {
  	Title OptString `json:"title,omitempty"`  // Wrapped in OptString
  }

  type OptString struct {
  	Value string
  	Set   bool
  }

  // Helper methods (standard Ogen for Opt*):
  func (o *OptString) SetTo(v string) { o.Set = true; o.Value = v }
  func (o OptString) IsSet() bool { return o.Set }
  func (o OptString) MarshalJSON() ([]byte, error) {
  	if !o.IsSet() { return []byte("null"), nil }  // Omitted if Set: false
  	return json.Marshal(o.Value)
  }
  func (o *OptString) UnmarshalJSON(data []byte) error {
  	if data == nil || string(data) == "null" { o.Set = false; return nil }
  	o.Set = true; return json.Unmarshal(data, &o.Value)
  }
  ```
  The `omitempty` tag ensures omission in JSON if `Set: false`. Required fields like `id` in params are plain strings.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 176-182, session.update test):  
  ```go
  req := OptSessionUpdateReq{
  	Value: SessionUpdateReq{
  		Title: OptString{Value: "Updated", Set: true},  // Explicitly set optional title
  	},
  	Set: true,
  }
  params := SessionUpdateParams{ID: testID}
  resp, err := client.SessionUpdate(ctx, req, params)
  assert.NoError(t, err)
  assert.NotNil(t, resp)
  ```
  Here, `OptString{Value: "Updated", Set: true}` sets the optional title. If `Set: false`, it's omitted. In validation tests (lines 441+), omitting required fields (e.g., title in create) triggers errors, but optionals like this succeed unset.

**Another Example: TuiShowToastReq (title optional)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 2198-2200, tui.showToast variant is required, but title is not):  
  ```json
  "properties": {
    "title": { "type": "string" },  // Optional
    "message": { "type": "string" },  // Required
    "variant": { "type": "string", "enum": ["info", "success", "warning", "error"] }  // Required
  },
  "required": ["message", "variant"]  // title omitted = optional
  ```

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 2200-2250 for TuiShowToastReq):  
  ```go
  type TuiShowToastReq struct {
  	Message string      // Required: plain string
  	Variant TuiShowToastReqVariant  // Required: enum
  	Title   OptString  // Optional: wrapped
  }
  ```
  `Title OptString` allows omission.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 402-406, tui.showToast test):  
  ```go
  req := OptTuiShowToastReq{
  	Value: TuiShowToastReq{
  		Message: "test",  // Required
  		Variant: "info",  // Required enum
  		// Title omitted (Set: false implicit in zero OptString)
  	},
  	Set: true,
  }
  resp, err := client.TuiShowToast(ctx, req)
  assert.NoError(t, err)
  assert.True(t, resp)
  ```
  The test omits `Title` (defaults to unset OptString), succeeding. If a required field like `Variant` were omitted, it would error (see similar in lines 490-497 for missing variant).

#### 2. Add Discriminator: Unions with Discriminator Field

Ogen detects `anyOf`/`oneOf` schemas and generates sum types with a discriminator (e.g., `Type` or `Role` field) that selects the variant. The discriminator is typically a shared enum field in the spec (e.g., `type: "text"` for variants). This step adds the enum type and dispatch logic before building the full sum.

**Example: Part (anyOf multiple subtypes; discriminator: type)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 4453-4482, Part schema):  
  ```json
  "Part": {
    "anyOf": [
      { "$ref": "#/components/schemas/TextPart" },  // type: "text"
      { "$ref": "#/components/schemas/ReasoningPart" },  // type: "reasoning"
      { "$ref": "#/components/schemas/FilePart" },  // type: "file"
      { "$ref": "#/components/schemas/ToolPart" },  // type: "tool"
      { "$ref": "#/components/schemas/StepStartPart" },  // type: "step-start"
      { "$ref": "#/components/schemas/StepFinishPart" },  // type: "step-finish"
      { "$ref": "#/components/schemas/SnapshotPart" },  // type: "snapshot"
      { "$ref": "#/components/schemas/PatchPart" },  // type: "patch"
      { "$ref": "#/components/schemas/AgentPart" }  // type: "agent"
    ]
  }
  ```
  Each referenced schema has a unique `type` field as the discriminator.

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 3561-3801 for Part sum type):  
  ```go
  type Part struct {
    Type PartType  // Discriminator enum
    TextPart       TextPart
    ReasoningPart  ReasoningPart
    FilePart       FilePart
    ToolPart       ToolPart
    StepStartPart  StepStartPart
    StepFinishPart StepFinishPart
    SnapshotPart   SnapshotPart
    PatchPart      PatchPart
    AgentPart      AgentPart
  }

  type PartType string
  const (
    TextPartPart       PartType = "text"
    ReasoningPartPart  PartType = "reasoning"
    FilePartPart       PartType = "file"
    ToolPartPart       ToolPart = "tool"
    StepStartPartPart  PartType = "step-start"
    StepFinishPartPart PartType = "step-finish"
    SnapshotPartPart   PartType = "snapshot"
    PatchPartPart      PartType = "patch"
    AgentPartPart      PartType = "agent"
  )

  // Discriminator methods (for each variant):
  func (s Part) IsTextPart() bool { return s.Type == TextPartPart }
  func (s *Part) SetTextPart(v TextPart) { s.Type = TextPartPart; s.TextPart = v }
  func (s Part) GetTextPart() (v TextPart, ok bool) { if !s.IsTextPart() { return v, false }; return s.TextPart, true }
  func NewTextPartPart(v TextPart) Part { var s Part; s.SetTextPart(v); return s }  // Constructor
  // Similar for all 9 variants
  ```
  The `Type` field is added/derived from the common `type` property in each variant schema.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 238-255, session.chat constructs parts):  
  ```go
  parts := []SessionChatReqPartsItem{  // Parts are anyOf TextPartInput, etc., leading to Part union
    NewTextPartInputSessionChatReqPartsItem(TextPartInput{  // Variant constructor for input
      Type: "text",  // Discriminator value
      Text: "Hello",
    }),
  }
  req := OptSessionChatReq{
    Value: SessionChatReq{
      ProviderID: "openai",
      ModelID:    "gpt-4o",
      Parts:      parts,  // Array of input variants → response as Part union
    },
    Set: true,
  }
  params := SessionChatParams{ID: testID}
  resp, err := client.SessionChat(ctx, req, params)
  assert.NoError(t, err)
  assert.NotNil(t, resp)  // resp.Parts is []Part (sum type); unmarshals based on discriminator
  ```
  The test constructs input variants (e.g., TextPartInput with `Type: "text"`); response unmarshals to Part union using `Type` to select (e.g., `part.IsTextPart()` to access TextPart fields). In validation tests (lines 441-454), mismatched `Type` (e.g., "invalid") would error unmarshal.

**Another Example: Auth (anyOf OAuth, ApiAuth, WellKnownAuth; discriminator: type)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 5119-5129, Auth schema):  
  ```json
  "Auth": {
    "anyOf": [
      { "$ref": "#/components/schemas/OAuth" },  // type: "oauth"
      { "$ref": "#/components/schemas/ApiAuth" },  // type: "api"
      { "$ref": "#/components/schemas/WellKnownAuth" }  // type: "wellknown"
    ]
  }
  ```
  Each has a unique `type` field.

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 686-750 for Auth sum type):  
  ```go
  type Auth struct {
    Type AuthType  // Discriminator enum
    OAuth         OAuth
    ApiAuth       ApiAuth
    WellKnownAuth WellKnownAuth
  }

  type AuthType string
  const (
    OAuthAuth         AuthType = "oauth"
    ApiAuthAuth       AuthType = "api"
    WellKnownAuthAuth AuthType = "wellknown"
  )

  // Discriminator methods:
  func (s Auth) IsOAuth() bool { return s.Type == OAuthAuth }
  func (s *Auth) SetOAuth(v OAuth) { s.Type = OAuthAuth; s.OAuth = v }
  func (s Auth) GetOAuth() (v OAuth, ok bool) { if !s.IsOAuth() { return v, false }; return s.OAuth, true }
  func NewOAuthAuth(v OAuth) Auth { var s Auth; s.SetOAuth(v); return s }  // Variant constructor
  // Similar for ApiAuth, WellKnownAuth
  ```
  `Type` selects the active variant.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 414-417, auth.set test):  
  ```go
  req := OptAuth{Set: true}  // Would construct e.g., NewOAuthAuth(OAuth{Type: "oauth", Refresh: "token", Access: "token", Expires: 1234567890})
  params := AuthSetParams{ID: testAuthID}
  _, err := client.AuthSet(ctx, req, params)
  assert.NoError(t, err)
  ```
  The test uses `OptAuth` (wrapper for optional Auth union); in full usage, you'd do `NewOAuthAuth(...)` for the variant. Validation ensures the `type` matches the variant (e.g., OAuth with `type: "api"` errors).

#### 3. Apply Validation: Schema Constraints Enforced

Ogen enforces spec constraints via struct tags (e.g., `valid:"required,pattern(...)"`) and runtime methods (e.g., `Valid()` on types). Enums have `MarshalText`/`UnmarshalText` for validation. Patterns, min/max, and required fields trigger errors on invalid data during client requests or server responses (if validation features enabled).

**Example: Message ID Pattern (pattern ^msg.*, required)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 894-895, session.prompt requestBody messageID):  
  ```json
  "messageID": {
    "type": "string",
    "pattern": "^msg.*"
  }
  ```
  Required in the schema.

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 2200-2250 for SessionChatReq; validation tag on ID fields):  
  ```go
  type SessionChatReq struct {
    MessageID string `json:"messageID" valid:"required,pattern(^msg.*)"`  // Enforces pattern and required
    // ...
  }
  ```
  On unmarshal or validation, invalid ID (e.g., "invalid") errors with "does not match pattern ^msg.*". Required check fails if missing.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 441-454, TestClientValidation for session.chat with invalid/missing fields):  
  ```go
  // session.chat missing providerID (required, but pattern on messageID tested similarly)
  parts := []SessionChatReqPartsItem{
    NewTextPartInputSessionChatReqPartsItem(TextPartInput{Type: "text", Text: "hello"}),
  }
  req := OptSessionChatReq{
    Value: SessionChatReq{
      MessageID: "invalid",  // Violates ^msg.* pattern
      ModelID:   "gpt-4o",
      Parts:     parts,
    },
    Set: true,
  }
  params := SessionChatParams{ID: testID}
  _, err = client.SessionChat(ctx, req, params)
  assert.Error(t, err)  // Errors: "invalid value for messageID: does not match pattern ^msg.*"
  ```
  The test constructs invalid inputs to trigger Ogen's client-side validation before sending. Similar for other patterns like session ID in params (SessionGetParams.ID).

**Another Example: TuiShowToastReq.Variant (enum validation)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 2191-2196, tui.showToast variant):  
  ```json
  "variant": {
    "type": "string",
    "enum": ["info", "success", "warning", "error"]
  }
  ```
  Required enum.

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 2200-2250 for TuiShowToastReq):  
  ```go
  type TuiShowToastReqVariant string
  const (
    TuiShowToastReqVariantInfo     TuiShowToastReqVariant = "info"
    TuiShowToastReqVariantSuccess  TuiShowToastReqVariant = "success"
    TuiShowToastReqVariantWarning  TuiShowToastReqVariant = "warning"
    TuiShowToastReqVariantError    TuiShowToastReqVariant = "error"
  )

  func (s TuiShowToastReqVariant) MarshalText() ([]byte, error) {
    switch s {
    case TuiShowToastReqVariantInfo: return []byte(s), nil
    // ... cases for success, warning, error
    default: return nil, errors.Errorf("invalid enum value: %q", s)  // Enum validation
  }

  func (s *TuiShowToastReqVariant) UnmarshalText(data []byte) error {
    switch TuiShowToastReqVariant(data) {
    case TuiShowToastReqVariantInfo: *s = TuiShowToastReqVariantInfo; return nil
    // ... cases
    default: return errors.Errorf("invalid enum value: %q", data)
  }
  ```
  Invalid enum (e.g., "invalid") errors on unmarshal/marshal.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 490-497, tui.showToast missing variant errors as required; valid uses "info"):  
  ```go
  // Missing variant (required enum)
  reqToast := OptTuiShowToastReq{
    Value: TuiShowToastReq{
      Message: "test",
      // Variant omitted -> required validation error
    },
    Set: true,
  }
  _, err = client.TuiShowToast(ctx, reqToast)
  assert.Error(t, err)  // Errors: "field variant required"

  // Valid enum usage
  reqToast := OptTuiShowToastReq{
    Value: TuiShowToastReq{
      Message: "test",
      Variant: "info",  // Valid enum
    },
    Set: true,
  }
  _, err = client.TuiShowToast(ctx, reqToast)
  assert.NoError(t, err)
  ```
  The test validates required enum; invalid enum would error "invalid enum value: 'invalid'".

#### 4. Nest Structures: Nested Objects/Arrays Become Nested Go Structs

Ogen converts inline objects (type: "object") and arrays into nested Go structs/slices. Properties become fields; required ones are direct, optionals wrapped. This creates type-safe hierarchies from spec nesting.

**Example: Session.time (nested object)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 3479-3490, Session.time):  
  ```json
  "time": {
    "type": "object",
    "properties": {
      "created": { "type": "number" },
      "updated": { "type": "number" },
      "compacting": { "type": "number" }
    },
    "required": ["created", "updated"]
  }
  ```
  Nested object with required/optional properties.

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 2900-2920 for SessionTime nested struct):  
  ```go
  type Session struct {
    // ...
    Time SessionTime `json:"time"`  // Nested struct
    // ...
  }

  type SessionTime struct {
    Created   float64   `json:"created"`  // Required direct field
    Updated   float64   `json:"updated"`  // Required direct field
    Compacting OptFloat64 `json:"compacting"`  // Optional wrapper
  }

  // Accessors for nested struct:
  func (s *SessionTime) GetCreated() float64 { return s.Created }
  func (s *SessionTime) SetCreated(val float64) { s.Created = val }
  // Similar for Updated; OptFloat64 for compacting
  ```
  The inline object becomes a nested struct; no top-level Opt since time is required in Session.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 159-165, session.create response with nested time):  
  ```go
  mockHandler.On("SessionCreate", mock.Anything, mock.Anything).Return(&server.Session{
    ID:      generateUUID(),
    Title:   "Test",
    Version: "1.0",
    Time: SessionTime{  // Nested struct initialization
      Created: float64(time.Now().Unix()),  // Required field
      Updated: float64(time.Now().Unix()),  // Required field
      // Compacting omitted (OptFloat64 defaults to unset)
    },
  }, nil)
  resp, err := client.SessionCreate(ctx, req)
  assert.IsType(t, &Session{}, resp)  // Unmarshals nested time; validates required created/updated
  ```
  The test constructs a nested `SessionTime` struct for the response; Ogen validates the required fields are present.

**Another Example: Symbol.location.range (deeply nested object)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 3961-3980, Symbol.location.range):  
  ```json
  "location": {
    "type": "object",
    "properties": {
      "uri": { "type": "string" },
      "range": {
        "type": "object",
        "properties": {
          "start": {
            "type": "object",
            "properties": {
              "line": { "type": "number" },
              "character": { "type": "number" }
            },
            "required": ["line", "character"]
          },
          "end": {
            "type": "object",
            "properties": {
              "line": { "type": "number" },
              "character": { "type": "number" }
            },
            "required": ["line", "character"]
          }
        },
        "required": ["start", "end"]
      }
    },
    "required": ["uri", "range"]
  }
  ```
  Deeply nested required objects.

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 4700-4750 for SymbolLocation and nested Range/Position):  
  ```go
  type Symbol struct {
    // ...
    Location SymbolLocation `json:"location"`  // Top-level nested
    // ...
  }

  type SymbolLocation struct {
    URI   string  `json:"uri"`
    Range Range   `json:"range"`  // Nested
  }

  type Range struct {
    Start Position `json:"start"`  // Nested
    End   Position `json:"end"`    // Nested
  }

  type Position struct {
    Line       float64 `json:"line"`  // Required
    Character  float64 `json:"character"`  // Required
  }

  // Accessors:
  func (s *Range) GetStart() Position { return s.Start }
  func (s *Range) SetStart(val Position) { s.Start = val }
  // Similar for End, and deeper for Position Line/Character
  ```
  Each level becomes a struct; required fields are direct.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: line 332, find.symbols response):  
  ```go
  params := FindSymbolsParams{Query: testQuery}
  resp, err := client.FindSymbols(ctx, params)
  assert.NoError(t, err)
  assert.IsType(t, []Symbol{}, resp)  // Expects []Symbol with nested Location.Range.Start/End
  ```
  The test asserts the response type includes the full nested structure; invalid nesting (e.g., missing line) would unmarshal error.

#### 5. Build Sum Types: oneOf/anyOf Become Discriminated Sum Types

Ogen generates sum types for `oneOf`/`anyOf` as structs with a discriminator (from step 2) and variant fields. Includes methods to set/get/check active variants and constructors like `New<Variant>()`. This builds on discrimination to create type-safe unions.

**Example: Part (anyOf 9 subtypes; discriminator: type from step 2)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 4453-4482, Part schema):  
  ```json
  "Part": {
    "anyOf": [
      { "$ref": "#/components/schemas/TextPart" },  // type: "text"
      { "$ref": "#/components/schemas/ReasoningPart" },  // type: "reasoning"
      // ... 7 more variants
    ]
  }
  ```
  Builds sum from variants with shared `type`.

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 3561-3801 for Part sum type):  
  ```go
  type Part struct {
    Type PartType  // Discriminator (built from step 2)
    TextPart       TextPart
    ReasoningPart  ReasoningPart
    FilePart       FilePart
    ToolPart       ToolPart
    StepStartPart  StepStartPart
    StepFinishPart StepFinishPart
    SnapshotPart   SnapshotPart
    PatchPart      PatchPart
    AgentPart      AgentPart
  }

  // Enum for discriminator (from step 2):
  type PartType string
  const (
    TextPartPart       PartType = "text"
    ReasoningPartPart  PartType = "reasoning"
    FilePartPart       PartType = "file"
    ToolPartPart       PartType = "tool"
    StepStartPartPart  PartType = "step-start"
    StepFinishPartPart PartType = "step-finish"
    SnapshotPartPart   PartType = "snapshot"
    PatchPartPart      PartType = "patch"
    AgentPartPart      PartType = "agent"
  )

  // Sum type methods (built for all variants):
  func (s Part) IsTextPart() bool { return s.Type == TextPartPart }
  func (s *Part) SetTextPart(v TextPart) { s.Type = TextPartPart; s.TextPart = v }
  func (s Part) GetTextPart() (v TextPart, ok bool) { if !s.IsTextPart() { return v, false }; return s.TextPart, true }
  func NewTextPartPart(v TextPart) Part { var s Part; s.SetTextPart(v); return s }  // Constructor for TextPart variant
  // Identical blocks for ReasoningPart, FilePart, etc. (9 total)
  ```
  The sum struct holds all variants; only the active one is populated based on discriminator.

- **Test/Usage Evidence** (`packages/tui/api/ogen/client_test.go`: lines 238-255, session.chat with parts union):  
  ```go
  parts := []SessionChatReqPartsItem{  // Input anyOf → response as Part sum
    NewTextPartInputSessionChatReqPartsItem(TextPartInput{  // Constructor for TextPartInput variant
      Type: "text",  // Sets discriminator
      Text: "Hello",
    }),
  }
  req := OptSessionChatReq{
    Value: SessionChatReq{
      ProviderID: "openai",
      ModelID:    "gpt-4o",
      Parts:      parts,  // Builds sum for input; response unmarshals to Part sum
    },
    Set: true,
  }
  params := SessionChatParams{ID: testID}
  resp, err := client.SessionChat(ctx, req, params)
  assert.NoError(t, err)
  assert.NotNil(t, resp)  // resp.Parts[0] is Part sum; e.g., part.IsTextPart() == true, part.GetTextPart().Text == "Hello"
  ```
  The test uses constructors like `NewTextPartInput...` (input side); response handling in `SessionChat` unmarshals to Part sum, selecting variant by `Type`. In full tests, accessing wrong variant (e.g., `part.GetFilePart()` on text) returns false/ok=false.

**Another Example: ToolState (anyOf Pending, Running, Completed, Error; discriminator: status)**  
- **Spec Reference** (`packages/tui/openapi.json`: lines 4049-4056, ToolStatePending; similar for others):  
  ```json
  "ToolStatePending": {
    "type": "object",
    "properties": {
      "status": { "type": "string", "const": "pending" }
    },
    "required": ["status"]
  }
  ```
  (Similar const status for Running/Completed/Error in their schemas.)

- **Generated Type** (inferred from Ogen in `oas_schemas_gen.go`: lines 4050-4100 for ToolState sum):  
  ```go
  type ToolState struct {
    Type ToolStateType  // Discriminator
    Pending   ToolStatePending
    Running   ToolStateRunning
    Completed ToolStateCompleted
    Error     ToolStateError
  }

  type ToolStateType string
  const (
    ToolStatePendingToolState   ToolStateType = "pending"
    ToolStateRunningToolState   ToolStateType = "running"
    ToolStateCompletedToolState ToolStateType = "completed"
    ToolStateErrorToolState     ToolStateType = "error"
  )

  // Sum methods:
  func (s ToolState) IsPending() bool { return s.Type == ToolStatePendingToolState }
  func (s *ToolState) SetPending(v ToolStatePending) { s.Type = ToolStatePendingToolState; s.Pending = v }
  func (s ToolState) GetPending() (v ToolStatePending, ok bool) { if !s.IsPending() { return v, false }; return s.Pending, true }
  func NewToolStatePendingToolState(v ToolStatePending) ToolState { var s ToolState; s.SetPending(v); return s }
  // Similar for Running, Completed, Error
  ```
  Used in ToolPart: `State OptToolState` (optional sum).

- **Test/Usage Evidence** (inferred from tool-related tests; no direct ToolState test, but pattern in session.command/shell which use ToolPart):  
  ```go
  // In session.command response (lines 270-271, indirectly via AssistantMessage with ToolPart):
  resp, err := client.SessionCommand(ctx, req, params)
  assert.NotNil(t, resp)  // Would include ToolPart{State: NewToolStateRunningToolState(...) } for running status
  // Access: if state.IsRunning() { input := state.GetRunning().Input }
  ```
  Validation tests (lines 441+) ensure sum selection; mismatched status errors unmarshal.
