package client

import (
	"context"

	opencode "github.com/sst/opencode-sdk-go"
)

// HealthClient defines the interface for health and metadata operations
type HealthClient interface {
	Ping(ctx context.Context) error
	Version(ctx context.Context) (*VersionInfo, error)
	Config(ctx context.Context) (*ServerConfig, error)
}

// VersionInfo represents version information
type VersionInfo struct {
	Version   string
	BuildDate string
	GitCommit string
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Providers []*Provider
}

// Provider represents a provider configuration
type Provider struct {
	ID          string
	Name        string
	Description *string
	Models      []*Model
}

// Model represents a model configuration
type Model struct {
	ID       string
	Name     string
	Provider *Provider
}

// LegacyHealth implements HealthClient using the legacy SDK
type LegacyHealth struct {
	client *opencode.Client
}

// NewLegacyHealth creates a new legacy health client
func NewLegacyHealth(client *opencode.Client) HealthClient {
	return &LegacyHealth{client: client}
}

func (l *LegacyHealth) Ping(ctx context.Context) error {
	// Use a lightweight config call as a health check
	_, err := l.client.Config.Get(ctx, opencode.ConfigGetParams{})
	return err
}

func (l *LegacyHealth) Version(ctx context.Context) (*VersionInfo, error) {
	// The legacy SDK does not expose a version endpoint; return empty info.
	return &VersionInfo{}, nil
}

func (l *LegacyHealth) Config(ctx context.Context) (*ServerConfig, error) {
	resp, err := l.client.App.Providers(ctx, opencode.AppProvidersParams{})
	if err != nil {
		return nil, err
	}
	cfg := &ServerConfig{Providers: []*Provider{}}
	if resp == nil {
		return cfg, nil
	}
	for _, p := range resp.Providers {
		prov := &Provider{ID: p.ID, Name: p.Name, Models: []*Model{}}
		// models are provided as a map[string]opencode.Model
		for mid, m := range p.Models {
			model := &Model{ID: mid, Name: m.Name, Provider: prov}
			prov.Models = append(prov.Models, model)
		}
		cfg.Providers = append(cfg.Providers, prov)
	}
	return cfg, nil
}
