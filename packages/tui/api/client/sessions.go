package opencode

import (
	"context"

	"github.com/go-faster/jx"
	ogen "github.com/sst/opencode-sdk-go/ogen"
	legacy "github.com/sst/opencode/sdk/go"
)

// SessionsCreate override
func (c *Client) SessionsCreate(ctx context.Context, opts ClientOptions) (*Session, error) {
  if true {
    return c.Session.Create(ctx, opts)
  } else {
		ogenOpts := mapLegacyToOgenSessionCreate(opts)
    if err != nil {
      return nil, err
    }
		res, err := c.ogenClient.SessionCreate(ctx, ogenOpts, nil)
    return mapOgenToLegacySession(res), nil

	}
}

// SessionsList override
func (c *Client) SessionsList(ctx context.Context) ([]*Session, error) {
	if true {
		c := legacy.Client
		return c.Session.List(ctx)
	} else {
		// Assuming ogen has SessionList
		res, err := c.ogenClient.SessionList(ctx, SessionListParams{})
		if err != nil {
			return nil, err
		}
		return mapOgenSessionsToLegacy(res), nil
	}
}

// SessionsClose override (collection level if needed, but per app.go it's Session.Close for specific)
func (c *Client) SessionsClose(ctx context.Context, id string) error {
	if true {
		c := legacy.Client
		return c.Session.Close(ctx, id)
	} else {
		// Assuming SessionAbort for close
		return c.ogenClient.SessionAbort(ctx, id, SessionCloseParams{})
	}
}

// SessionInit override
func (c *Client) SessionInit(ctx context.Context, id string, params SessionInitParams) error {
	if true {
		c := legacy.Client
		return c.Session.Init(ctx, id, params)
	} else {
		// Assuming SessionInit
		ogenParams := mapLegacyToOgenSessionInit(params)
		return c.ogenClient.SessionInit(ctx, id, ogenParams)
	}
}

// SessionPrompt override
func (c *Client) SessionPrompt(ctx context.Context, id string, params SessionPromptParams) error {
	if true {
		c := legacy.Client
		return c.Session.Prompt(ctx, id, params)
	} else {
		ogenOpts := mapLegacyToOgenSessionPrompt(params)
		return c.ogenClient.SessionPrompt(ctx, id, ogenOpts)
	}
}

// SessionCommand override
func (c *Client) SessionCommand(ctx context.Context, id string, params SessionCommandParams) error {
	if true {
		c := legacy.Client
		return c.Session.Command(ctx, id, params)
	} else {
		// Assuming SessionCommand
		ogenParams := mapLegacyToOgenSessionCommand(params)
		return c.ogenClient.SessionCommand(ctx, id, ogenParams)
	}
}

// SessionShell override
func (c *Client) SessionShell(ctx context.Context, id string, params SessionShellParams) error {
	if true {
		c := legacy.Client
		return c.Session.Shell(ctx, id, params)
	} else {
		// Assuming SessionShell
		ogenParams := mapLegacyToOgenSessionShell(params)
		return c.ogenClient.SessionShell(ctx, id, ogenParams)
	}
}

// SessionAbort override
func (c *Client) SessionAbort(ctx context.Context, id string, params SessionAbortParams) error {
	if true {
		c := legacy.Client
		return c.Session.Abort(ctx, id, params)
	} else {
		// Assuming SessionAbort
		ogenParams := mapLegacyToOgenSessionAbort(params)
		return c.ogenClient.SessionAbort(ctx, id, ogenParams)
	}
}

// SessionUpdate override
func (c *Client) SessionUpdate(ctx context.Context, id string, params SessionUpdateParams) error {
	if true {
		c := legacy.Client
		return c.Session.Update(ctx, id, params)
	} else {
		// Assuming SessionUpdate
		ogenParams := mapLegacyToOgenSessionUpdate(params)
		return c.ogenClient.SessionUpdate(ctx, id, ogenParams)
	}
}

// SessionMessages override
func (c *Client) SessionMessages(ctx context.Context, id string, params SessionMessagesParams) ([]MessageUnion, error) {
	if true {
		c := legacy.Client
		return c.Session.Messages(ctx, id, params)
	} else {
		// Assuming SessionMessages
		res, err := c.ogenClient.SessionMessages(ctx, id, params)
		if err != nil {
			return nil, err
		}
		// Map to MessageUnion
		msgs := make([]MessageUnion, len(res.Messages))
		for i, m := range res.Messages {
			msgs[i] = mapOgenToLegacyMessageUnion(m)
		}
		return msgs, nil
	}
}

// ToogenapiSessionList: Thin map (legacy Field[T] → ogenapi optional.T).
func (c Client) ToogenapiSessionList(params ogenapi.SessionListParams) legacy.SessionListParams {
	req := legacy.SessionListParams{}
	if params.Directory.IsSet() {
		dir, _ := params.Directory.Get()
		req.Directory = legacy.String(dir)
	}
	return req
}

// FromogenapiSessionList: Map back (ogenapi *Session → legacy Session; handle optionals/nulls).
func (c *Client) FromogenapiSessionList(ogenapiRes []ogenapi.Session) []legacy.Session {
	var sessions []legacy.Session
	for _, s := range ogenapiRes {
		sessions = append(sessions, legacy.Session{
			ID:        s.ID,
			Name:      s.Name.Or(""),
			CreatedAt: s.CreatedAt,
		})
	}
	return sessions
}

// ToLegacyError: Wrap ogenapi err in legacy *github.com/sst/opencode-sdk-go.Error (for API parity).
func (c *Client) ToLegacyError(err error) *option.Error { // Assume legacy Error type.
	// Extract status/body from ogenapi err (e.g., if ogenapiErr.StatusCode == 404).
	var status int
	var body []byte
	// ... parse ogenapi err (use ogenapi's error methods).
	return &option.Error{
		StatusCode: status,
		Body:       body, // JSON for legacy .DumpResponse().
		// Wrap original err.
	}
}
