# Shared Response Types

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go/shared">shared</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go/shared#MessageAbortedError">MessageAbortedError</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go/shared">shared</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go/shared#ProviderAuthError">ProviderAuthError</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go/shared">shared</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go/shared#UnknownError">UnknownError</a>

# Event

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#EventListResponse">EventListResponse</a>

Methods:

- <code title="get /event">client.Event.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#EventService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#EventListResponse">EventListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# App

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Agent">Agent</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#App">App</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Model">Model</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Provider">Provider</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppProvidersResponse">AppProvidersResponse</a>

Methods:

- <code title="get /agent">client.App.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppService.Agents">Agents</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Agent">Agent</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /app">client.App.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#App">App</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /app/init">client.App.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppService.Init">Init</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /log">client.App.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppService.Log">Log</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppLogParams">AppLogParams</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /config/providers">client.App.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppService.Providers">Providers</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AppProvidersResponse">AppProvidersResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Find

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Symbol">Symbol</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindTextResponse">FindTextResponse</a>

Methods:

- <code title="get /find/file">client.Find.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindService.Files">Files</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindFilesParams">FindFilesParams</a>) ([]<a href="https://pkg.go.dev/builtin#string">string</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /find/symbol">client.Find.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindService.Symbols">Symbols</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindSymbolsParams">FindSymbolsParams</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Symbol">Symbol</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /find">client.Find.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindService.Text">Text</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindTextParams">FindTextParams</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FindTextResponse">FindTextResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# File

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#File">File</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FileReadResponse">FileReadResponse</a>

Methods:

- <code title="get /file">client.File.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FileService.Read">Read</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, query <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FileReadParams">FileReadParams</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FileReadResponse">FileReadResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /file/status">client.File.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FileService.Status">Status</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#File">File</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Config

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Config">Config</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#KeybindsConfig">KeybindsConfig</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#McpLocalConfig">McpLocalConfig</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#McpRemoteConfig">McpRemoteConfig</a>

Methods:

- <code title="get /config">client.Config.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#ConfigService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Config">Config</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Command

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Command">Command</a>

Methods:

- <code title="get /command">client.Command.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#CommandService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Command">Command</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Session

Params Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AgentPartInputParam">AgentPartInputParam</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FilePartInputParam">FilePartInputParam</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FilePartSourceUnionParam">FilePartSourceUnionParam</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FilePartSourceTextParam">FilePartSourceTextParam</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FileSourceParam">FileSourceParam</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SymbolSourceParam">SymbolSourceParam</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TextPartInputParam">TextPartInputParam</a>

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AgentPart">AgentPart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AssistantMessage">AssistantMessage</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FilePart">FilePart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FilePartSource">FilePartSource</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FilePartSourceText">FilePartSourceText</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#FileSource">FileSource</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Message">Message</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Part">Part</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#ReasoningPart">ReasoningPart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SnapshotPart">SnapshotPart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#StepFinishPart">StepFinishPart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#StepStartPart">StepStartPart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SymbolSource">SymbolSource</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TextPart">TextPart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#ToolPart">ToolPart</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#ToolStateCompleted">ToolStateCompleted</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#ToolStateError">ToolStateError</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#ToolStatePending">ToolStatePending</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#ToolStateRunning">ToolStateRunning</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#UserMessage">UserMessage</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionChatResponse">SessionChatResponse</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionCommandResponse">SessionCommandResponse</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionMessageResponse">SessionMessageResponse</a>
- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionMessagesResponse">SessionMessagesResponse</a>

Methods:

- <code title="post /session">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionNewParams">SessionNewParams</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="patch /session/{id}">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Update">Update</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionUpdateParams">SessionUpdateParams</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /session">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /session/{id}">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Delete">Delete</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/abort">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Abort">Abort</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/message">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Chat">Chat</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionChatParams">SessionChatParams</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionChatResponse">SessionChatResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /session/{id}/children">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Children">Children</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/command">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Command">Command</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionCommandParams">SessionCommandParams</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionCommandResponse">SessionCommandResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /session/{id}">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/init">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Init">Init</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionInitParams">SessionInitParams</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /session/{id}/message/{messageID}">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Message">Message</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, messageID <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionMessageResponse">SessionMessageResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /session/{id}/message">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Messages">Messages</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) ([]<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionMessagesResponse">SessionMessagesResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/revert">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Revert">Revert</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionRevertParams">SessionRevertParams</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/share">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Share">Share</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/shell">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Shell">Shell</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionShellParams">SessionShellParams</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#AssistantMessage">AssistantMessage</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/summarize">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Summarize">Summarize</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionSummarizeParams">SessionSummarizeParams</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /session/{id}/unrevert">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Unrevert">Unrevert</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="delete /session/{id}/share">client.Session.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionService.Unshare">Unshare</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>) (<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Session">Session</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Permissions

Response Types:

- <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#Permission">Permission</a>

Methods:

- <code title="post /session/{id}/permissions/{permissionID}">client.Session.Permissions.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionPermissionService.Respond">Respond</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, id <a href="https://pkg.go.dev/builtin#string">string</a>, permissionID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#SessionPermissionRespondParams">SessionPermissionRespondParams</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Tui

Methods:

- <code title="post /tui/append-prompt">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.AppendPrompt">AppendPrompt</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiAppendPromptParams">TuiAppendPromptParams</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/clear-prompt">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.ClearPrompt">ClearPrompt</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/execute-command">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.ExecuteCommand">ExecuteCommand</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiExecuteCommandParams">TuiExecuteCommandParams</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/open-help">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.OpenHelp">OpenHelp</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/open-models">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.OpenModels">OpenModels</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/open-sessions">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.OpenSessions">OpenSessions</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/open-themes">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.OpenThemes">OpenThemes</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/show-toast">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.ShowToast">ShowToast</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, body <a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go">opencode</a>.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiShowToastParams">TuiShowToastParams</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /tui/submit-prompt">client.Tui.<a href="https://pkg.go.dev/git.j9xym.com/opencode-api-go#TuiService.SubmitPrompt">SubmitPrompt</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (<a href="https://pkg.go.dev/builtin#bool">bool</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
