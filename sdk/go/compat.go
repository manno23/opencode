package opencode

import (
	"context"

	"github.com/sst/opencode-sdk-go/option"
)

// NewCompatClient creates a new API client with the given base URL
func NewCompatClient(baseURL string) *Client {
	return NewClient(option.WithBaseURL(baseURL))
}

// SessionInit initializes a new session (placeholder for future implementation)
func SessionInit(ctx context.Context, c *Client, params map[string]any) (any, error) {
	// This is a stub — update with real mapping as API stabilizes
	return nil, nil
}
