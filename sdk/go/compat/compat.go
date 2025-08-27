package compat

import (
	"context"
	opencodesdk "git.j9xym.com/opencode-api-go"
)

// Minimal compatibility adapters used by TUI to call common operations

func NewClient(baseURL string) (*opencodesdk.Client, error) {
	return opencodesdk.NewClient(baseURL)
}

func SessionInit(ctx context.Context, c *opencodesdk.Client, params map[string]any) (any, error) {
	// adapt to new SDK SessionInit equivalent
	// This is a stub — update with real mapping as API stabilizes
	return nil, nil
}
