package client

import (
	"context"
	"fmt"
	"io"

	"github.com/sst/opencode-sdk-go"
)

// MessagesClient defines the interface for message operations
type MessagesClient interface {
	// Core operations
	Send(ctx context.Context, sessionID, content string, opts SendOptions) (*Message, error)
	Get(ctx context.Context, messageID string) (*Message, error)
	List(ctx context.Context, sessionID string, opts MessageListOptions) ([]*Message, error)
	
	// Message actions
	Revert(ctx context.Context, messageID string) error
	Unrevert(ctx context.Context, messageID string) error
	
	// Streaming
	Stream(ctx context.Context, sessionID, content string, opts StreamOptions) (<-chan *StreamEvent, error)
}

// Message represents a message in the facade layer
type Message struct {
	ID        string
	SessionID string
	Content   string
	Role      MessageRole
	CreatedAt string
	UpdatedAt string
	Parts     []*MessagePart
}

// MessageRole represents message roles
type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)

// MessagePart represents different types of message content
type MessagePart struct {
	Type    MessagePartType
	Content string
	Meta    map[string]interface{}
}

// MessagePartType represents different message part types
type MessagePartType string

const (
	MessagePartTypeText MessagePartType = "text"
	MessagePartTypeTool MessagePartType = "tool"
	MessagePartTypeFile MessagePartType = "file"
)

// Request/Response types
type SendOptions struct {
	Attachments []io.Reader
	Stream      bool
}

type MessageListOptions struct {
	Limit  *int
	Offset *int
}

type StreamOptions struct {
	Attachments []io.Reader
}

type StreamEvent struct {
	Type MessagePartType
	Data interface{}
	Done bool
}

// LegacyMessages implements MessagesClient using the legacy SDK
type LegacyMessages struct {
	client *opencode.Client
}

// NewLegacyMessages creates a new legacy messages client
func NewLegacyMessages(client *opencode.Client) MessagesClient {
	return &LegacyMessages{client: client}
}

func (l *LegacyMessages) Send(ctx context.Context, sessionID, content string, opts SendOptions) (*Message, error) {
	return nil, fmt.Errorf("legacy SDK does not support facade Messages.Send")
}

func (l *LegacyMessages) Get(ctx context.Context, messageID string) (*Message, error) {
	return nil, fmt.Errorf("legacy SDK does not support facade Messages.Get")
}

func (l *LegacyMessages) List(ctx context.Context, sessionID string, opts MessageListOptions) ([]*Message, error) {
	return []*Message{}, nil
}

func (l *LegacyMessages) Revert(ctx context.Context, messageID string) error {
	return fmt.Errorf("legacy SDK does not support facade Messages.Revert")
}

func (l *LegacyMessages) Unrevert(ctx context.Context, messageID string) error {
	return fmt.Errorf("legacy SDK does not support facade Messages.Unrevert")
}

func (l *LegacyMessages) Stream(ctx context.Context, sessionID, content string, opts StreamOptions) (<-chan *StreamEvent, error) {
	return nil, fmt.Errorf("legacy SDK does not support facade Messages.Stream")
}

