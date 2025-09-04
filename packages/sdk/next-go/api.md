# OpenCode SDK API Stubs

This file contains stubbed type definitions for all symbols referenced in the TUI project that are missing from the opencode package.

## Types

### Core Types

```go
type Client struct {
    // TODO: Implement client
}

type Session struct {
    ID       string
    ParentID string
    Title    string
    Share    *SessionShare
    Revert   *SessionRevert
}

type Permission struct {
    ID        string
    SessionID string
    Metadata  map[string]any
    MessageID string
    CallID    string
}

type MessageUnion interface {
    // Union type for messages
}

type PartUnion interface {
    // Union type for parts
}

type Project struct {
    // TODO: Implement project
}

type Agent struct {
    // TODO: Implement agent
}

type Provider struct {
    // TODO: Implement provider
}

type Model struct {
    // TODO: Implement model
}

type Config struct {
    Share ConfigShare
}

type Command struct {
    // TODO: Implement command
}

type Path struct {
    // TODO: Implement path
}

type Symbol struct {
    // TODO: Implement symbol
}
```

### Message Types

```go
type UserMessage struct {
    Role UserMessageRole
    Time UserMessageTime
    Parts []PartUnion
}

type AssistantMessage struct {
    Parts []PartUnion
    Error *AssistantMessageError
}

type AssistantMessageError struct {
    Message AssistantMessageErrorMessage
}

type AssistantMessageErrorMessage struct {
    OutputLengthError *AssistantMessageErrorMessageOutputLengthError
}

type AssistantMessageErrorMessageOutputLengthError struct {
    // TODO: Implement
}
```

### Part Types

```go
type TextPart struct {
    Type TextPartType
    Text string
    Time *TextPartTime
}

type FilePart struct {
    Type FilePartType
    Source FilePartSource
}

type AgentPart struct {
    Source AgentPartSource
}

type ToolPart struct {
    State ToolPartState
}

type ReasoningPart struct {
    // TODO: Implement reasoning part
}

type StepStartPart struct {
    // TODO: Implement step start part
}

type StepFinishPart struct {
    // TODO: Implement step finish part
}
```

### Parameter Types

```go
type AppLogParams struct {
    Service *string
    Level   *AppLogParamsLevel
    Message *string
    Extra   *map[string]any
}

type SessionGetParams struct {
    // TODO: Add fields
}

type SessionChildrenParams struct {
    // TODO: Add fields
}

type SessionShareParams struct {
    // TODO: Add fields
}

type SessionUnshareParams struct {
    // TODO: Add fields
}

type SessionRevertParams struct {
    MessageID *string
}

type SessionUnrevertParams struct {
    // TODO: Add fields
}

type SessionNewParams struct {
    // TODO: Add fields
}

type SessionPromptParams struct {
    Model *SessionPromptParamsModel
    PartUnion []SessionPromptParamsPartUnion
}

type SessionInitParams struct {
    // TODO: Add fields
}

type SessionSummarizeParams struct {
    // TODO: Add fields
}

type SessionCommandParams struct {
    // TODO: Add fields
}

type SessionShellParams struct {
    // TODO: Add fields
}

type SessionAbortParams struct {
    // TODO: Add fields
}

type SessionListParams struct {
    // TODO: Add fields
}

type SessionDeleteParams struct {
    // TODO: Add fields
}

type SessionUpdateParams struct {
    // TODO: Add fields
}

type SessionMessagesParams struct {
    // TODO: Add fields
}

type SessionMessageParams struct {
    // TODO: Add fields
}

type ConfigGetParams struct {
    // TODO: Add fields
}

type CommandListParams struct {
    // TODO: Add fields
}

type AppProvidersParams struct {
    // TODO: Add fields
}

type AgentListParams struct {
    // TODO: Add fields
}

type FindSymbolsParams struct {
    // TODO: Add fields
}

type FindFilesParams struct {
    // TODO: Add fields
}

type FileStatusParams struct {
    // TODO: Add fields
}

type PathGetParams struct {
    // TODO: Add fields
}

type ProjectCurrentParams struct {
    // TODO: Add fields
}

type EventListParams struct {
    // TODO: Add fields
}
```

### Response Types

```go
type AppProvidersResponse struct {
    Providers []Provider
}

type SessionPermissionRespondParams struct {
    Response *int
}
```

### Source and Input Types

```go
type FilePartSource struct {
    Type FilePartSourceType
    Text *FilePartSourceText
    File *FilePartSourceFile
    Symbol *FilePartSourceSymbol
}

type FilePartSourceText struct {
    Text string
}

type FilePartSourceFile struct {
    // TODO: Implement
}

type FilePartSourceSymbol struct {
    Range SymbolSourceRange
}

type AgentPartSource struct {
    // TODO: Implement
}

type FileSourceParam struct {
    Type FileSourceType
}

type SymbolSourceParam struct {
    Range SymbolSourceRange
}

type TextPartInputParam struct {
    Type TextPartInputType
    Text string
    Time *TextPartInputTime
}

type FilePartInputParam struct {
    Type FilePartInputType
}

type AgentPartInputParam struct {
    Type AgentPartInputType
    Source AgentPartInputSource
}

type AgentPartInputSourceParam struct {
    // TODO: Implement
}

type FilePartSourceUnionParam struct {
    // TODO: Implement
}

type SymbolSourceRangeParam struct {
    Start SymbolSourceRangeStart
    End SymbolSourceRangeEnd
}

type SymbolSourceRangeStartParam struct {
    // TODO: Implement
}

type SymbolSourceRangeEndParam struct {
    // TODO: Implement
}
```

### State and Error Types

```go
type ToolPartState struct {
    Status ToolPartStateStatus
}

type MessageAbortedError struct {
    // TODO: Implement
}
```

### Session Related Types

```go
type SessionShare struct {
    URL string
}

type SessionRevert struct {
    MessageID string
    PartID string
    Diff *string
}
```

## Constants

```go
const (
    AppLogParamsLevelDebug AppLogParamsLevel = iota
    AppLogParamsLevelInfo
    AppLogParamsLevelWarn
    AppLogParamsLevelError
)

const (
    ConfigShareDisabled ConfigShare = iota
)

const (
    UserMessageRoleUser UserMessageRole = "user"
)

const (
    TextPartTypeText TextPartType = "text"
)

const (
    FilePartTypeFile FilePartType = "file"
)

const (
    TextPartInputTypeText TextPartInputType = "text"
)

const (
    FilePartInputTypeFile FilePartInputType = "file"
)

const (
    AgentPartInputTypeAgent AgentPartInputType = "agent"
)

const (
    FilePartSourceTypeFile FilePartSourceType = "file"
    FilePartSourceTypeSymbol FilePartSourceType = "symbol"
)

const (
    FileSourceTypeFile FileSourceType = "file"
)

const (
    SymbolSourceTypeSymbol SymbolSourceType = "symbol"
)

const (
    ToolPartStateStatusPending ToolPartStateStatus = "pending"
    ToolPartStateStatusCompleted ToolPartStateStatus = "completed"
    ToolPartStateStatusError ToolPartStateStatus = "error"
)

const (
    AgentModePrimary AgentMode = "primary"
)

const (
    SessionPermissionRespondParamsResponseOnce int = 0
    SessionPermissionRespondParamsResponseAlways int = 1
    SessionPermissionRespondParamsResponseReject int = 2
)
```

## Event Types

```go
type EventListResponseEventInstallationUpdated struct{}
type EventListResponseEventIdeInstalled struct{}
type EventListResponseEventSessionDeleted struct{}
type EventListResponseEventSessionUpdated struct{}
type EventListResponseEventSessionError struct{}
type EventListResponseEventMessageUpdated struct{}
type EventListResponseEventMessageRemoved struct{}
type EventListResponseEventMessagePartUpdated struct{}
type EventListResponseEventMessagePartRemoved struct{}
type EventListResponseEventPermissionUpdated struct{}
type EventListResponseEventPermissionReplied struct{}
```

## Error Types

```go
type ProviderAuthError struct{}
type UnknownError struct{}
```

## Functions

```go
func NewClient(serverURL string, opts ...ClientOption) (*Client, error) {
    // TODO: Implement
    return &Client{}, nil
}

func F[T any](v T) *T {
    return &v
}
```

## Type Aliases

```go
type UserMessageRole string
type UserMessageTime string
type TextPartTime string
type TextPartInputTimeParam string
type AppLogParamsLevel int
type ConfigShare int
type TextPartType string
type FilePartType string
type TextPartInputType string
type FilePartInputType string
type AgentPartInputType string
type FilePartSourceType string
type FileSourceType string
type SymbolSourceType string
type ToolPartStateStatus string
type AgentMode string
type SessionPromptParamsModel string
type SessionPromptParamsPartUnion interface{}
```

## Range Types

```go
type SymbolSourceRange struct {
    Start SymbolSourceRangeStart
    End SymbolSourceRangeEnd
}

type SymbolSourceRangeStart struct {
    // TODO: Implement
}

type SymbolSourceRangeEnd struct {
    // TODO: Implement
}
```

## Client Options

```go
type ClientOption interface {
    // TODO: Implement
}
```

## TODO

This is a comprehensive list of stub definitions. Each TODO indicates areas that need proper implementation based on the OpenAPI specification.
