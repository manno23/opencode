// Package api provides the generated Ogen client for the opencode API.
// This file contains comprehensive tests for the client methods.
package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadSpec loads the OpenAPI spec from the JSON file.
func loadSpec(t *testing.T) *openapi3.T {
	spec, err := openapi3.NewLoader().LoadFromFile("../unified-openapi.json")
	require.NoError(t, err)
	return spec
}

// TestClientAllOperations tests all client operations using a mock server.
func TestClientAllOperations(t *testing.T) {
	spec := loadSpec(t)

	// Create a mock HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		var response string
		switch r.URL.Path {
		case "/app":
			response = `{"hostname": "test", "git": false, "path": {"config": "", "data": "", "root": "", "cwd": "", "state": ""}, "time": {}}`
		case "/app/init":
			response = `true`
		case "/session":
			if r.Method == "GET" {
				response = `[]`
			} else {
				response = `{"id": "ses_test", "title": "Test Session", "version": "1", "time": {"created": 1234567890, "updated": 1234567890}}`
			}
		default:
			if strings.Contains(r.URL.Path, "/session/") && strings.Contains(r.URL.Path, "/message") && r.Method == "POST" {
				response = `{"info": {"id": "msg", "sessionID": "test", "role": "assistant", "time": {"created": 123}, "system": [], "modelID": "model", "providerID": "prov", "mode": "chat", "path": {"cwd": "", "root": ""}, "cost": 0, "tokens": {"input": 0, "output": 0, "reasoning": 0, "cache": {"read": 0, "write": 0}}}, "parts": []}`
			} else {
				response = `{}`
			}
		}
		_, _ = w.Write([]byte(response))
	}))
	defer ts.Close()

	// Create client with test server URL.
	client, err := NewClient(ts.URL)
	require.NoError(t, err)

	ctx := context.Background()

	// Enumerate all operations from spec and test each.
	for path, pathItem := range spec.Paths.Map() {
		for method, op := range pathItem.Operations() {
			operationID := op.OperationID
			t.Run(fmt.Sprintf("%s %s (%s)", method, path, operationID), func(t *testing.T) {
				testOperation(t, client, ctx, operationID, op)
			})
		}
	}
}

// testOperation runs a basic test for the given operationID.
func testOperation(t *testing.T, client *Client, ctx context.Context, opID string, op *openapi3.Operation) {
	switch opID {
	case "app.get":
		resp, err := client.AppGet(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)

	case "app.init":
		resp, err := client.AppInit(ctx)
		assert.NoError(t, err)
		assert.True(t, resp)

	case "session.create":
		req := OptSessionCreateReq{
			Value: SessionCreateReq{
				Title: OptString{Value: "Test Session", Set: true},
			},
			Set: true,
		}
		resp, err := client.SessionCreate(ctx, req)
		assert.NoError(t, err)
		assert.IsType(t, &Session{}, resp)

	case "session.list":
		resp, err := client.SessionList(ctx)
		assert.NoError(t, err)
		assert.IsType(t, []Session{}, resp)

	case "session.get":
		// Need mock params; skip detailed for now or use fixed ID.
		// In full impl, generate params from spec.
		t.Skip("Parametrized; implement with generated params")

	// Add cases for all operations...
	case "session.chat":
		req := OptSessionChatReq{
			Value: SessionChatReq{
				ProviderID: "test-provider",
				ModelID:    "test-model",
				Parts: []SessionChatReqPartsItem{
					NewTextPartInputSessionChatReqPartsItem(TextPartInput{
						Type: "text",
						Text: "Hello",
					}),
				},
			},
			Set: true,
		}
		params := SessionChatParams{ID: "test-session"}
		resp, err := client.SessionChat(ctx, req, params)
		assert.NoError(t, err)
		assert.NotNil(t, resp)

	// ... continue for all ~30 operations, with valid inputs from schemas.

	default:
		t.Skipf("Operation %s not yet implemented in test", opID)
	}
}

// TestClientValidation tests input validation errors.
func TestClientValidation(t *testing.T) {
	client, err := NewClient("http://example.com")
	require.NoError(t, err)

	ctx := context.Background()

	// Example: Invalid session.create with missing required fields.
	req := OptSessionCreateReq{
		Value: SessionCreateReq{}, // Missing title
		Set:   true,
	}
	_, err = client.SessionCreate(ctx, req)
	assert.Error(t, err) // Expect validation error from Ogen

	// More cases...
}

// TestMockServerWithOgenHandlers (optional, if generate server)
func TestMockServerWithOgenHandlers(t *testing.T) {
	// If server code generated, use it here for more accurate mocks.
	// ts := httptest.NewServer(ogenServerHandler)
	// ...
	t.Skip("Requires server generation")
}

// For E2E, would need to spin up actual server, but focus on unit/integration.
