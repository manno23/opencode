package opencode

import (
	"context"
	"fmt"

	"github.com/sst/opencode-api-go/internal/requestconfig"
)

// SessionServiceInterface defines the unified interface for session operations
type SessionServiceInterface interface {
	Get(ctx context.Context, id string, params SessionGetParams) (Session, error)
	New(ctx context.Context, params SessionNewParams) (*Session, error)
	Update(ctx context.Context, id string, params SessionUpdateParams) (*Session, error)
	List(ctx context.Context, query SessionListParams) (*[]Session, error)
	Delete(ctx context.Context, id string, body SessionDeleteParams) (*bool, error)
}

// GeneratedSessionAdapter implements SessionServiceInterface using the generated API client
// This adapter handles parameter conversion and response mapping between
// the generated API types and custom types
type GeneratedSessionAdapter struct {
	client *Client
}

func NewGeneratedSessionAdapter(client *Client) *GeneratedSessionAdapter {
	return &GeneratedSessionAdapter{client: client}
}

func (a *GeneratedSessionAdapter) Get(ctx context.Context, id string, params SessionGetParams) (Session, error) {
	if id == "" {
		return Session{}, fmt.Errorf("missing required id parameter")
	}

	// Convert custom params to API params
	directory := ""
	if params.Directory.Present {
		directory = params.Directory.Value
	}

	// Call generated API (placeholder)
	// apiSession, err := a.client.apiClient.SessionGet(ctx, api.SessionGetParams{ID: id, Directory: directory})
	// if err != nil {
	//     return Session{}, fmt.Errorf("failed to fetch session from API: %w", err)
	// }

	// For demonstration, return a placeholder session
	// In real implementation, this would map from api.Session to Session
	return Session{
		ID:        id,
		Directory: directory,
		ProjectID: "project-123",
		Version:   "1.0",
		Time: SessionTime{
			Created: 1234567890,
			Updated: 1234567890,
		},
	}, nil
}

// Extend the pattern to other methods
func (a *GeneratedSessionAdapter) New(ctx context.Context, params SessionNewParams) (*Session, error) {
	// Convert params and call generated API
	// apiParams := a.convertNewParams(params)
	// apiSession, err := a.client.apiClient.SessionNew(ctx, apiParams)
	// return a.mapToCustomSession(apiSession), err

	// Placeholder implementation
	return &Session{
		ID:        "new-session-id",
		Directory: "/tmp",
		ProjectID: "project-123",
		Version:   "1.0",
	}, nil
}

func (a *GeneratedSessionAdapter) Update(ctx context.Context, id string, params SessionUpdateParams) (*Session, error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	// Convert and call generated API
	// apiParams := a.convertUpdateParams(id, params)
	// apiSession, err := a.client.apiClient.SessionUpdate(ctx, apiParams)
	// return a.mapToCustomSession(apiSession), err

	return &Session{
		ID:        id,
		Directory: "/updated/path",
		ProjectID: "project-123",
		Version:   "1.1",
	}, nil
}

func (a *GeneratedSessionAdapter) List(ctx context.Context, query SessionListParams) (*[]Session, error) {
	// Convert query params and call generated API
	// apiParams := a.convertListParams(query)
	// apiSessions, err := a.client.apiClient.SessionList(ctx, apiParams)
	// return a.mapToCustomSessions(apiSessions), err

	sessions := []Session{
		{ID: "session-1", Directory: "/path1", ProjectID: "project-123"},
		{ID: "session-2", Directory: "/path2", ProjectID: "project-456"},
	}
	return &sessions, nil
}

func (a *GeneratedSessionAdapter) Delete(ctx context.Context, id string, body SessionDeleteParams) (*bool, error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	// Call generated API
	// apiParams := a.convertDeleteParams(id, body)
	// result, err := a.client.apiClient.SessionDelete(ctx, apiParams)
	// return &result, err

	result := true
	return &result, nil
}

// CustomSessionAdapter implements SessionServiceInterface using custom request handling
// This adapter uses the existing ExecuteNewRequest pattern
type CustomSessionAdapter struct {
	client *Client
}

func NewCustomSessionAdapter(client *Client) *CustomSessionAdapter {
	return &CustomSessionAdapter{client: client}
}

func (a *CustomSessionAdapter) Get(ctx context.Context, id string, params SessionGetParams) (Session, error) {
	var res Session
	path := fmt.Sprintf("session/%s", id)
	err := requestconfig.ExecuteNewRequest(ctx, "GET", path, params, &res)
	return res, err
}

func (a *CustomSessionAdapter) New(ctx context.Context, params SessionNewParams) (*Session, error) {
	var res Session
	err := requestconfig.ExecuteNewRequest(ctx, "POST", "session", params, &res)
	return &res, err
}

func (a *CustomSessionAdapter) Update(ctx context.Context, id string, params SessionUpdateParams) (*Session, error) {
	var res Session
	path := fmt.Sprintf("session/%s", id)
	err := requestconfig.ExecuteNewRequest(ctx, "PATCH", path, params, &res)
	return &res, err
}

func (a *CustomSessionAdapter) List(ctx context.Context, query SessionListParams) (*[]Session, error) {
	var res []Session
	err := requestconfig.ExecuteNewRequest(ctx, "GET", "session", query, &res)
	return &res, err
}

func (a *CustomSessionAdapter) Delete(ctx context.Context, id string, body SessionDeleteParams) (*bool, error) {
	var res bool
	path := fmt.Sprintf("session/%s", id)
	err := requestconfig.ExecuteNewRequest(ctx, "DELETE", path, body, &res)
	return &res, err
}

// UnifiedSessionAdapter combines both approaches, preferring generated but falling back to custom
// This provides a seamless interface that can use either backend implementation
type UnifiedSessionAdapter struct {
	generated *GeneratedSessionAdapter
	custom    *CustomSessionAdapter
}

func NewUnifiedSessionAdapter(client *Client) *UnifiedSessionAdapter {
	return &UnifiedSessionAdapter{
		generated: NewGeneratedSessionAdapter(client),
		custom:    NewCustomSessionAdapter(client),
	}
}

func (a *UnifiedSessionAdapter) Get(ctx context.Context, id string, params SessionGetParams) (Session, error) {
	// Try generated API first for better performance/features
	session, err := a.generated.Get(ctx, id, params)
	if err != nil {
		// Fallback to custom implementation
		return a.custom.Get(ctx, id, params)
	}
	return session, nil
}

func (a *UnifiedSessionAdapter) New(ctx context.Context, params SessionNewParams) (*Session, error) {
	session, err := a.generated.New(ctx, params)
	if err != nil {
		return a.custom.New(ctx, params)
	}
	return session, nil
}

func (a *UnifiedSessionAdapter) Update(ctx context.Context, id string, params SessionUpdateParams) (*Session, error) {
	session, err := a.generated.Update(ctx, id, params)
	if err != nil {
		return a.custom.Update(ctx, id, params)
	}
	return session, nil
}

func (a *UnifiedSessionAdapter) List(ctx context.Context, query SessionListParams) (*[]Session, error) {
	sessions, err := a.generated.List(ctx, query)
	if err != nil {
		return a.custom.List(ctx, query)
	}
	return sessions, nil
}

func (a *UnifiedSessionAdapter) Delete(ctx context.Context, id string, body SessionDeleteParams) (*bool, error) {
	result, err := a.generated.Delete(ctx, id, body)
	if err != nil {
		return a.custom.Delete(ctx, id, body)
	}
	return result, nil
}

// SessionServiceAdapterFactory creates the appropriate adapter based on configuration
type SessionServiceAdapterFactory struct {
	preferGenerated bool
}

func NewSessionServiceAdapterFactory(preferGenerated bool) *SessionServiceAdapterFactory {
	return &SessionServiceAdapterFactory{preferGenerated: preferGenerated}
}

func (f *SessionServiceAdapterFactory) CreateAdapter(client *Client) SessionServiceInterface {
	if f.preferGenerated {
		return NewUnifiedSessionAdapter(client)
	}
	return NewCustomSessionAdapter(client)
}

// Usage examples:
//
// 1. Using unified adapter (recommended):
// factory := NewSessionServiceAdapterFactory(true)
// service := factory.CreateAdapter(client)
// session, err := service.Get(ctx, "session-id", params)
//
// 2. Using specific adapters:
// generated := NewGeneratedSessionAdapter(client)
// custom := NewCustomSessionAdapter(client)
//
// 3. Direct unified usage:
// adapter := NewUnifiedSessionAdapter(client)
// session, err := adapter.Get(ctx, "session-id", params)
