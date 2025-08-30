package opencode

import (
	"context"

	"git.j9xym.com/opencode-api-go"
	"git.j9xym.com/opencode-api-go/option"
)

// NewClient creates a new API client with the given base URL
func NewClient(baseURL string) (*opencode.Client, error) {
	return opencode.NewClient(option.WithBaseURL(baseURL)), nil
}

// SessionInit initializes a new session (placeholder for future implementation)
func SessionInit(ctx context.Context, c *opencode.Client, params map[string]any) (any, error) {
	// This is a stub — update with real mapping as API stabilizes
	return nil, nil
}
