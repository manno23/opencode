package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	opencode "github.com/sst/opencode-sdk-go"
	api "github.com/sst/opencode/api/ogen"
)

// OgenSessions implements SessionsClient using the generated ogen client
type OgenSessions struct {
	client *api.Client
}

// NewOgenSessions creates a new OgenSessions client
func NewOgenSessions(config *ClientConfig) SessionsClient {
	serverURL := os.Getenv("OPENCODE_SERVER")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	// Create ogen client with custom HTTP client for auth
	httpClient := &http.Client{}

	client, err := api.NewClient(serverURL,
		api.WithClient(httpClient),
	)
	if err != nil {
		// Fallback to legacy if ogen client creation fails
		return NewLegacySessions(config.LegacyClient)
	}

	return &OgenSessions{client: client}
}

// Create creates a new session using ogen client
func (o *OgenSessions) Create(ctx context.Context, opts CreateOpts) (*Session, error) {
	// Build ogen request from facade options
	req := api.OptSessionCreateReq{}
	if opts.Title != "" {
		req.Set = true
		req.Value = api.SessionCreateReq{
			Title: api.OptString{Value: opts.Title, Set: true},
		}
		if opts.ParentID != nil {
			req.Value.ParentID = api.OptString{Value: *opts.ParentID, Set: true}
		}
	}

	// Call ogen client
	resp, err := o.client.SessionCreate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ogen session create failed: %w", err)
	}

	// Map ogen response to facade type
	return mapOgenSessionCreateResToFacade(resp), nil
}

// List lists all sessions using ogen client
func (o *OgenSessions) List(ctx context.Context) ([]*Session, error) {
	// Call ogen client
	sessions, err := o.client.SessionList(ctx)
	if err != nil {
		return nil, fmt.Errorf("ogen session list failed: %w", err)
	}

	// Map ogen sessions to facade sessions
	result := make([]*Session, len(sessions))
	for i, s := range sessions {
		result[i] = mapOgenSessionToFacade(&s)
	}
	return result, nil
}

// Get gets a session by ID using ogen client
func (o *OgenSessions) Get(ctx context.Context, id string) (*Session, error) {
	params := api.SessionGetParams{
		ID: id,
	}

	session, err := o.client.SessionGet(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ogen session get failed: %w", err)
	}

	return mapOgenSessionToFacade(session), nil
}

// Delete deletes a session by ID using ogen client
func (o *OgenSessions) Delete(ctx context.Context, id string) error {
	params := api.SessionDeleteParams{
		ID: id,
	}

	_, err := o.client.SessionDelete(ctx, params)
	if err != nil {
		return fmt.Errorf("ogen session delete failed: %w", err)
	}
	return nil
}

// mapOgenSessionToFacade converts ogen Session to facade Session
func mapOgenSessionToFacade(ogenSession *api.Session) *Session {
	if ogenSession == nil {
		return nil
	}

	// Map ParentID from OptString to string (opencode.Session expects string, not *string)
	parentID := ""
	if ogenSession.ParentID.Set {
		parentID = ogenSession.ParentID.Value
	}

	// Create a new Session using the opencode.Session type alias
	// We need to create the full struct since Session = opencode.Session
	return &Session{
		ID:       ogenSession.ID,
		Title:    ogenSession.Title,
		ParentID: parentID,
		Version:  ogenSession.Version,
		// Note: Directory, ProjectID, and other required fields would need to be mapped
		// from the ogen Session struct. For now, using defaults or empty values.
		Directory: "", // Would need to be mapped from ogen Session
		ProjectID: "", // Would need to be mapped from ogen Session
		Time: opencode.SessionTime{
			Created: ogenSession.Time.Created,
			Updated: ogenSession.Time.Updated,
		},
		// Other fields like Revert, Share would need proper mapping
	}
}

// mapOgenSessionCreateResToFacade converts SessionCreateRes to facade Session
func mapOgenSessionCreateResToFacade(resp api.SessionCreateRes) *Session {
	// SessionCreateRes can be either *Session or *Error
	// We need to check which type it is and extract the Session

	// Try to assert as Session pointer
	if session, ok := resp.(*api.Session); ok {
		return mapOgenSessionToFacade(session)
	}

	// If we can't extract a session, return nil
	return nil
}

// formatTimestamp converts float64 timestamp to string
func formatTimestamp(timestamp float64) string {
	// Convert Unix timestamp to time
	t := time.Unix(int64(timestamp), 0)
	return t.Format(time.RFC3339)
}
