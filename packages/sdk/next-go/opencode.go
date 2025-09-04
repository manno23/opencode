package opencode

import (
	"context"
	"encoding/json"
	apiPkg "github.com/sst/opencode-sdk-go/api"
)

// F helper
func F[T any](v T) *T { return &v }

// Client wraps the generated API client
type Client struct {
	apiClient *apiPkg.Client
	App       *AppService
	Session   *SessionService
}

// NewClient creates a new Client with the generated API client
func NewClient(serverURL string, opts ...apiPkg.ClientOption) (*Client, error) {
	apiClient, err := apiPkg.NewClient(serverURL, opts...)
	if err != nil {
		return nil, err
	}
	client := &Client{apiClient: apiClient}
	client.App = &AppService{client: client}
	client.Session = &SessionService{
		client:      client,
		Permissions: &SessionPermissionsService{client: client},
	}
	return client, nil
}

func (c *Client) Get(ctx context.Context, path string, query any, v any) error { return nil }
func (c *Client) Post(ctx context.Context, path string, body any, v any) error { return nil }

// App service
type AppService struct {
	client *Client
}

type AppLogParamsLevel int

const (
	AppLogParamsLevelDebug AppLogParamsLevel = iota
	AppLogParamsLevelInfo
	AppLogParamsLevelWarn
	AppLogParamsLevelError
)

type AppLogParams struct {
	Service *string
	Level   *AppLogParamsLevel
	Message *string
	Extra   *map[string]any
}

func (s *AppService) Log(ctx context.Context, params AppLogParams) (any, error) {
	// Convert to generated types
	var level apiPkg.AppLogReqLevel
	switch *params.Level {
	case AppLogParamsLevelDebug:
		level = apiPkg.AppLogReqLevelDebug
	case AppLogParamsLevelInfo:
		level = apiPkg.AppLogReqLevelInfo
	case AppLogParamsLevelWarn:
		level = apiPkg.AppLogReqLevelWarn
	case AppLogParamsLevelError:
		level = apiPkg.AppLogReqLevelError
	}

	req := apiPkg.AppLogReq{
		Service: *params.Service,
		Level:   level,
		Message: *params.Message,
	}
	if params.Extra != nil {
		// Convert map[string]any to map[string][]byte
		extraMap := make(apiPkg.AppLogReqExtra)
		for k, v := range *params.Extra {
			// Marshal v to JSON bytes
			data, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			extraMap[k] = data
		}
		req.Extra = apiPkg.OptAppLogReqExtra{Value: extraMap, Set: true}
	}

	optReq := apiPkg.OptAppLogReq{Value: req, Set: true}
	apiParams := apiPkg.AppLogParams{} // No query params for this endpoint

	return s.client.apiClient.AppLog(ctx, optReq, apiParams)
}

// Session service
type SessionService struct {
	client      *Client
	Permissions *SessionPermissionsService
}

type SessionPermissionsService struct {
	client *Client
}

func (s *SessionService) Get(ctx context.Context, id string, params any) (Session, error) {
	apiParams := apiPkg.SessionGetParams{ID: id}
	apiSession, err := s.client.apiClient.SessionGet(ctx, apiParams)
	if err != nil {
		return Session{}, err
	}
	var parentID string
	if apiSession.ParentID.Set {
		parentID = apiSession.ParentID.Value
	}
	var share *SessionShare
	if apiSession.Share.Set {
		share = &SessionShare{} // Populate fields as needed
	}
	var revert *SessionRevert
	if apiSession.Revert.Set {
		revert = &SessionRevert{} // Populate fields as needed
	}
	return Session{
		ID:       apiSession.ID,
		ParentID: parentID,
		Title:    apiSession.Title,
		Share:    share,
		Revert:   revert,
	}, nil
}

func (s *SessionPermissionsService) Respond(ctx context.Context, sessionID string, permissionID string, body SessionPermissionRespondParams) (any, error) {
	var response apiPkg.PostSessionByIdPermissionsByPermissionIDReqResponse
	switch *body.Response {
	case SessionPermissionRespondParamsResponseOnce:
		response = apiPkg.PostSessionByIdPermissionsByPermissionIDReqResponseOnce
	case SessionPermissionRespondParamsResponseAlways:
		response = apiPkg.PostSessionByIdPermissionsByPermissionIDReqResponseAlways
	case SessionPermissionRespondParamsResponseReject:
		response = apiPkg.PostSessionByIdPermissionsByPermissionIDReqResponseReject
	}

	apiParams := apiPkg.PostSessionByIdPermissionsByPermissionIDParams{
		ID:           sessionID,
		PermissionID: permissionID,
	}
	apiBody := apiPkg.PostSessionByIdPermissionsByPermissionIDReq{
		Response: response,
	}

	return s.client.apiClient.PostSessionByIdPermissionsByPermissionID(ctx, apiPkg.OptPostSessionByIdPermissionsByPermissionIDReq{Value: apiBody, Set: true}, apiParams)
}

// Session types
type Session struct {
	ID       string
	ParentID string
	Title    string
	Share    *SessionShare
	Revert   *SessionRevert
}

type SessionShare struct {
	// Add fields as needed
}

type SessionRevert struct {
	// Add fields as needed
}

type SessionPermissionRespondParams struct {
	Response *int
}

const (
	SessionPermissionRespondParamsResponseOnce   = 0
	SessionPermissionRespondParamsResponseAlways = 1
	SessionPermissionRespondParamsResponseReject = 2
)

// Permission
type Permission struct {
	ID        string
	SessionID string
	Metadata  *map[string]any
	MessageID string
	CallID    string
}

// AppLog response placeholder
type AppLogResponse struct{}

