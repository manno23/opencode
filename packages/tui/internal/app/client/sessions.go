package client

import (
	"context"

	opencode "github.com/sst/opencode-sdk-go"
)

// SessionsClient defines the interface for session operations
type SessionsClient interface {
	Create(ctx context.Context, opts CreateOpts) (*Session, error)
	List(ctx context.Context) ([]*Session, error)
	Get(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
}

// Session represents a session, matching the legacy SDK shape
type Session = opencode.Session

// CreateOpts contains options for creating a session
type CreateOpts struct {
	ParentID *string
	Title    string
}
