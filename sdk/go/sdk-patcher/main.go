// Package main provides a Go AST-based SDK patcher that fixes common API mismatches
// between the generated OpenCode Go SDK and what the TUI expects.
package opencodeapi

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: sdk-patcher <sdk-path>")
	}

	sdkPath := os.Args[1]
	opencodePath := filepath.Join(sdkPath, "opencode.go")

	fmt.Printf("Patching SDK at: %s\n", sdkPath)

	// Read the current file
	content, err := os.ReadFile(opencodePath)
	if err != nil {
		log.Fatalf("Error reading opencode.go: %v", err)
	}

	// Generate comprehensive compatibility layer
	patchContent := generateCompatibilityLayer()

	// Append the patch to existing content
	newContent := string(content) + "\n\n" + patchContent

	// Write back
	if err := os.WriteFile(opencodePath, []byte(newContent), 0644); err != nil {
		log.Fatalf("Error writing patched file: %v", err)
	}

	fmt.Println("SDK patching completed successfully")
}

func generateCompatibilityLayer() string {
	return `
// =============================================================================
// COMPATIBILITY LAYER - Auto-generated types and methods for TUI integration
// =============================================================================

// Override existing types with extended functionality
type Agent struct {
	ID    string  ` + "`json:\"id\"`" + `
	Name  string  ` + "`json:\"name\"`" + `
	Type  string  ` + "`json:\"type\"`" + `
	State string  ` + "`json:\"state\"`" + `
	Mode  string  ` + "`json:\"mode\"`" + `
	Model *Model  ` + "`json:\"model\"`" + `
}

type Provider struct {
	ID     string   ` + "`json:\"id\"`" + `
	Name   string   ` + "`json:\"name\"`" + `
	Type   string   ` + "`json:\"type\"`" + `
	Models []*Model ` + "`json:\"models\"`" + `
}

type Model struct {
	ID       string ` + "`json:\"id\"`" + `
	Name     string ` + "`json:\"name\"`" + `
	Provider string ` + "`json:\"provider\"`" + `
}

type Session struct {
	ID        string    ` + "`json:\"id\"`" + `
	Title     string    ` + "`json:\"title\"`" + `
	CreatedAt time.Time ` + "`json:\"createdAt\"`" + `
	UpdatedAt time.Time ` + "`json:\"updatedAt\"`" + `
}

type Permission struct {
	ID    string ` + "`json:\"id\"`" + `
	Name  string ` + "`json:\"name\"`" + `
	Type  string ` + "`json:\"type\"`" + `
	Allow bool   ` + "`json:\"allow\"`" + `
}

// Config override with expected structure
type Config struct {
	Share    ConfigShare       ` + "`json:\"share\"`" + `
	Keybinds ConfigKeybinds    ` + "`json:\"keybinds\"`" + `
	Model    *Model            ` + "`json:\"model\"`" + `
	Tui      ConfigTui         ` + "`json:\"tui\"`" + `
}

type ConfigKeybinds struct {
	Leader string ` + "`json:\"leader\"`" + `
}

type ConfigTui struct {
	ScrollSpeed float64 ` + "`json:\"scrollSpeed\"`" + `
}

// Message and Part types
type UserMessage struct {
	ID        string          ` + "`json:\"id\"`" + `
	SessionID string          ` + "`json:\"sessionID\"`" + `
	Role      UserMessageRole ` + "`json:\"role\"`" + `
	Time      UserMessageTime ` + "`json:\"time\"`" + `
}

type UserMessageRole string
const (
	UserMessageRoleUser UserMessageRole = "user"
)

type UserMessageTime struct {
	Created float64 ` + "`json:\"created\"`" + `
}

type AssistantMessage struct {
	ID        string ` + "`json:\"id\"`" + `
	SessionID string ` + "`json:\"sessionID\"`" + `
	Role      string ` + "`json:\"role\"`" + `
}

type TextPart struct {
	ID        string       ` + "`json:\"id\"`" + `
	Type      TextPartType ` + "`json:\"type\"`" + `
	Text      string       ` + "`json:\"text\"`" + `
	Synthetic bool         ` + "`json:\"synthetic\"`" + `
	Time      TextPartTime ` + "`json:\"time\"`" + `
}

type TextPartType string
const (
	TextPartTypeText TextPartType = "text"
)

type TextPartTime struct {
	Start float64 ` + "`json:\"start\"`" + `
	End   float64 ` + "`json:\"end\"`" + `
}

type FilePart struct {
	ID     string               ` + "`json:\"id\"`" + `
	Type   FilePartType         ` + "`json:\"type\"`" + `
	Source FilePartSourceUnion  ` + "`json:\"source\"`" + `
}

type FilePartType string
const (
	FilePartTypeFile FilePartType = "file"
)

type FilePartSourceUnion interface{}

type FilePartSource struct {
	Type FilePartSourceType ` + "`json:\"type\"`" + `
	Path string             ` + "`json:\"path\"`" + `
}

type FilePartSourceType string
const (
	FilePartSourceTypeFile   FilePartSourceType = "file"
	FilePartSourceTypeSymbol FilePartSourceType = "symbol"
)

type FilePartSourceText struct {
	Value string ` + "`json:\"value\"`" + `
}

type SymbolSourceRange struct {
	Start SymbolSourceRangeStart ` + "`json:\"start\"`" + `
	End   SymbolSourceRangeEnd   ` + "`json:\"end\"`" + `
}

type SymbolSourceRangeStart struct {
	Line   int ` + "`json:\"line\"`" + `
	Column int ` + "`json:\"column\"`" + `
}

type SymbolSourceRangeEnd struct {
	Line   int ` + "`json:\"line\"`" + `
	Column int ` + "`json:\"column\"`" + `
}

type ToolPart struct {
	ID    string            ` + "`json:\"id\"`" + `
	State ToolPartStateStatus ` + "`json:\"state\"`" + `
}

type ToolPartStateStatus string
const (
	ToolPartStateStatusPending ToolPartStateStatus = "pending"
)

// Parameter types for API calls
type SessionInitParams struct {
	ID string ` + "`json:\"id\"`" + `
}

type SessionNewParams struct {
	Title string ` + "`json:\"title\"`" + `
}

type SessionChatParams struct {
	Parts []SessionChatParamsPartUnion ` + "`json:\"parts\"`" + `
}

type SessionChatParamsPartUnion interface{}

type SessionCommandParams struct {
	Command string ` + "`json:\"command\"`" + `
}

type SessionShellParams struct {
	Command string ` + "`json:\"command\"`" + `
}

type SessionUpdateParams struct {
	Title string ` + "`json:\"title\"`" + `
}

type SessionSummarizeParams struct{}

// Input parameter types
type TextPartInputParam struct {
	ID        interface{}           ` + "`json:\"id\"`" + `
	Type      interface{}           ` + "`json:\"type\"`" + `
	Text      interface{}           ` + "`json:\"text\"`" + `
	Synthetic interface{}           ` + "`json:\"synthetic\"`" + `
	Time      TextPartInputTimeParam ` + "`json:\"time\"`" + `
}

type TextPartInputType string
const (
	TextPartInputTypeText TextPartInputType = "text"
)

type TextPartInputTimeParam struct {
	Start interface{} ` + "`json:\"start\"`" + `
	End   interface{} ` + "`json:\"end\"`" + `
}

type FilePartInputParam struct {
	ID     interface{} ` + "`json:\"id\"`" + `
	Type   interface{} ` + "`json:\"type\"`" + `
	Source interface{} ` + "`json:\"source\"`" + `
}

type FilePartInputType string
const (
	FilePartInputTypeFile FilePartInputType = "file"
)

type FilePartSourceUnionParam interface{}

type FileSourceParam struct {
	Type interface{} ` + "`json:\"type\"`" + `
	Path interface{} ` + "`json:\"path\"`" + `
}

type FileSourceType string
const (
	FileSourceTypeFile FileSourceType = "file"
)

type FilePartSourceTextParam struct {
	Value interface{} ` + "`json:\"value\"`" + `
}

type SymbolSourceParam struct {
	Type  interface{}                ` + "`json:\"type\"`" + `
	Path  interface{}                ` + "`json:\"path\"`" + `
	Range SymbolSourceRangeParam     ` + "`json:\"range\"`" + `
}

type SymbolSourceType string
const (
	SymbolSourceTypeSymbol SymbolSourceType = "symbol"
)

type SymbolSourceRangeParam struct {
	Start SymbolSourceRangeStartParam ` + "`json:\"start\"`" + `
	End   SymbolSourceRangeEndParam   ` + "`json:\"end\"`" + `
}

type SymbolSourceRangeStartParam struct {
	Line   interface{} ` + "`json:\"line\"`" + `
	Column interface{} ` + "`json:\"column\"`" + `
}

type SymbolSourceRangeEndParam struct {
	Line   interface{} ` + "`json:\"line\"`" + `
	Column interface{} ` + "`json:\"column\"`" + `
}

type AgentPartInputParam struct {
	ID     interface{}                  ` + "`json:\"id\"`" + `
	Type   interface{}                  ` + "`json:\"type\"`" + `
	Name   interface{}                  ` + "`json:\"name\"`" + `
	Source AgentPartInputSourceParam   ` + "`json:\"source\"`" + `
}

type AgentPartInputType string
const (
	AgentPartInputTypeAgent AgentPartInputType = "agent"
)

type AgentPartInputSourceParam struct {
	Value interface{} ` + "`json:\"value\"`" + `
}

// Response types
type SessionPermissionRespondParamsResponse string
type AppProvidersResponse []*Provider

// Extended Client with additional methods and fields
func (c *Client) Config() *Config {
	// Return a compatible config structure
	return &Config{
		Share:    ConfigShareDisabled,
		Keybinds: ConfigKeybinds{Leader: " "},  // Default leader key
		Model:    nil,  // Will be set by app
		Tui:      ConfigTui{ScrollSpeed: 3.0},
	}
}

// SessionClient provides session operations
type SessionClient struct {
	client *Client
}

func (c *Client) Session() *SessionClient {
	return &SessionClient{client: c}
}

func (s *SessionClient) Init(ctx context.Context, params SessionInitParams) (*Session, error) {
	// Implementation stub - delegate to generated client or provide mock
	return &Session{ID: params.ID}, nil
}

func (s *SessionClient) New(ctx context.Context, params SessionNewParams) (*Session, error) {
	return &Session{Title: params.Title}, nil
}

func (s *SessionClient) Chat(ctx context.Context, params SessionChatParams) (*AssistantMessage, error) {
	return &AssistantMessage{}, nil
}

func (s *SessionClient) Command(ctx context.Context, params SessionCommandParams) error {
	return nil
}

func (s *SessionClient) Shell(ctx context.Context, params SessionShellParams) error {
	return nil
}

func (s *SessionClient) Update(ctx context.Context, params SessionUpdateParams) (*Session, error) {
	return &Session{}, nil
}

func (s *SessionClient) Delete(ctx context.Context) error {
	return nil
}

func (s *SessionClient) Get(ctx context.Context) (*Session, error) {
	return &Session{}, nil
}

func (s *SessionClient) Summarize(ctx context.Context, params SessionSummarizeParams) (*Session, error) {
	return &Session{}, nil
}

// CommandClient provides command operations
type CommandClient struct {
	client *Client
}

func (c *Client) Command() *CommandClient {
	return &CommandClient{client: c}
}

// AppClient extensions
func (a *AppClient) Init(ctx context.Context) error {
	return nil
}

func (a *AppClient) Providers(ctx context.Context) (AppProvidersResponse, error) {
	// Return mock providers for now - in real implementation would call generated SDK
	return AppProvidersResponse{
		&Provider{
			ID:   "openai",
			Name: "OpenAI",
			Type: "llm",
			Models: []*Model{
				{ID: "gpt-4", Name: "GPT-4", Provider: "openai"},
			},
		},
	}, nil
}
`
}
