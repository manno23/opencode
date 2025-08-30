package api

import (
	"context"
	"encoding/json"

	"git.j9xym.com/opencode-api-go"
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Request struct {
	Path string          `json:"path"`
	Body json.RawMessage `json:"body"`
}

func Start(ctx context.Context, program *tea.Program, client *opencode.Client) {
	// TODO: Implement TUI control API or remove if not needed
	// The /tui/control/next and /tui/control/response endpoints don't exist in the current API
}

func Reply(ctx context.Context, client *opencode.Client, response interface{}) tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement TUI control response API or remove if not needed
		return nil
	}
}
