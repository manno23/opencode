package client

import (
	opencode "github.com/sst/opencode-sdk-go"
)

// Client is the facade client that embeds the legacy client and provides facades for different domains
type Client struct {
	*opencode.Client
	Sessions SessionsClient
}

// NewClient creates a new Client with the legacy client embedded and default legacy implementations
func NewClient(legacyClient *opencode.Client) *Client {
	return &Client{
		Client:   legacyClient,
		Sessions: NewLegacySessions(legacyClient),
	}
}
