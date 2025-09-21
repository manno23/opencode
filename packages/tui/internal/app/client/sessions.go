package client

import (
	"context"

	opencode "github.com/sst/opencode-sdk-go"
)

// Session is a facade alias to the SDK session type.
// This avoids duplicating types and keeps the facade light-weight.
type Session = opencode.Session

// CreateOpts contains options to create a new session.
type CreateOpts struct {
	ParentID *string
	Title    string
}

// SessionsClient defines the interface for session operations
// used by the TUI app layer.
type SessionsClient interface {
	Create(ctx context.Context, opts CreateOpts) (*Session, error)
	List(ctx context.Context) ([]*Session, error)
	Get(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
}
