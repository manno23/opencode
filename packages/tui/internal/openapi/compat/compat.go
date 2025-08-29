package opencode

import (
	"context"
)

// NewClient creates a new API client with the given base URL
func NewClient(baseURL string) (*Client, error) {
	return NewClient(baseURL)
}

// SessionInit initializes a new session (placeholder for future implementation)
func SessionInit(ctx context.Context, c *Client, params map[string]any) (any, error) {
	// This is a stub — update with real mapping as API stabilizes
	return nil, nil
}
