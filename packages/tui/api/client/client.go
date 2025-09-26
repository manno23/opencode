package opencode

import (
	"net/http"
	"os"

	ogenapi "github.com/sst/opencode-sdk-go/ogen" // Local legacy (via replace)
	legacy "github.com/sst/opencode/sdk/go"       // Local legacy (via replace)
	"github.com/sst/opencode/sdk/go/option"
)

// ClientOptions captures important data similar to legacy SDK's options
type ClientOptions struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Client struct with embedded legacy for passthrough and ogen for overrides
type Client struct {
	Client 				 *legacy.Client // embedded for passthrough
	ogenClient     *ogenapi.Client
}

// NewClient initializes both legacy and ogenapi clients
func NewClient(opts ...option.RequestOption) (*Client, error) {
	// Append default options like legacy
	opts = append(opts, legacy.DefaultClientOptions()...)

	// Create ClientOptions from env vars
	clientOpts := ClientOptions{
		BaseURL:    os.Getenv("OPENCODE_BASE_URL"),
		HTTPClient: http.DefaultClient,
	}
	if clientOpts.BaseURL == "" {
		clientOpts.BaseURL = "https://api.opencode.dev" // default URL
	}

	// Initialize legacy client
	legacyClient := legacy.NewClient(opts...)

	// Initialize ogenapi client with ClientOptions
	var ogenClient *ogenapi.Client
	if false {  // Ogen
		ogenOpts := []ogenapi.ClientOption{
			ogenapi.WithClient(clientOpts.HTTPClient),
		}
		var err error
		ogenClient, err = ogenapi.NewClient(clientOpts.BaseURL, ogenOpts...)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		Client:     legacyClient,
		ogenClient: ogenClient,
	}, nil
}
