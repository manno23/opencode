package client

import (
	"net/http"
	"time"

	"github.com/sst/opencode-api-go"
	"github.com/sst/opencode-api-go/api"
)

// Client wraps the OpenCode API client with enhanced configuration
type Client struct {
	*opencode.Client
	baseURL   string
	authToken string
	Command   *opencode.CommandService
}

// Config holds the client configuration
type Config struct {
	BaseURL    string
	AuthToken  string
	Timeout    time.Duration
	MaxRetries int
}

// NewClient creates a new client with proper HTTP configuration
func NewClient(config Config) (*Client, error) {
	// Create the retry round tripper with exponential backoff
	retryTransport := &retryRoundTripper{
		maxRetries: config.MaxRetries,
		base:       http.DefaultTransport,
	}

	// Add authentication if token provided
	var transport http.RoundTripper = retryTransport
	if config.AuthToken != "" {
		transport = &authRoundTripper{
			token: config.AuthToken,
			base:  retryTransport,
		}
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	// Create SDK client with custom HTTP client
	sdkClient, err := opencode.NewClient(config.BaseURL, api.WithClient(httpClient))
	if err != nil {
		return nil, err
	}

	client := &Client{
		Client:    sdkClient,
		baseURL:   config.BaseURL,
		authToken: config.AuthToken,
	}
	client.Command = opencode.NewCommandService()
	return client, nil
}
