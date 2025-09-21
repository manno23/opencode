package client

import (
	"context"

	opencode "github.com/sst/opencode-sdk-go"
)

// LegacySessions implements SessionsClient using the legacy SDK
type LegacySessions struct {
	client *opencode.Client
}

// NewLegacySessions creates a new LegacySessions
func NewLegacySessions(client *opencode.Client) SessionsClient {
	return &LegacySessions{client: client}
}

// Create creates a new session
func (l *LegacySessions) Create(ctx context.Context, opts CreateOpts) (*Session, error) {
	params := opencode.SessionNewParams{}
	if opts.ParentID != nil {
		params.ParentID = opencode.F(*opts.ParentID)
	}
	if opts.Title != "" {
		params.Title = opencode.F(opts.Title)
	}
	result, err := l.client.Session.New(ctx, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List lists all sessions
func (l *LegacySessions) List(ctx context.Context) ([]*Session, error) {
	res, err := l.client.Session.List(ctx, opencode.SessionListParams{})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return []*Session{}, nil
	}
	sessions := make([]*Session, len(*res))
	for i, s := range *res {
		session := s // Copy to avoid pointer issues
		sessions[i] = &session
	}
	return sessions, nil
}

// Get gets a session by ID
func (l *LegacySessions) Get(ctx context.Context, id string) (*Session, error) {
	return l.client.Session.Get(ctx, id, opencode.SessionGetParams{})
}

// Delete deletes a session by ID
func (l *LegacySessions) Delete(ctx context.Context, id string) error {
	_, err := l.client.Session.Delete(ctx, id, opencode.SessionDeleteParams{})
	return err
}
