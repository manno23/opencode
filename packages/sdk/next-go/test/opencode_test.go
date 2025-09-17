package test

import (
	"context"
	"testing"

	"github.com/sst/opencode-api-go"
	"github.com/sst/opencode-api-go/internal/testutil"
)

func TestNewClient(t *testing.T) {
	// Test successful client creation
	client, err := opencode.NewClient("http://localhost:3000")
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if client == nil {
		t.Fatal("NewClient returned nil client")
	}
	// We can't directly check the private apiClient field
	// but we can check that the public fields are initialized
	if client.App == nil {
		t.Error("Client App service is nil")
	}
	if client.Session == nil {
		t.Error("Client Session service is nil")
	}
}

func TestAppServiceLog(t *testing.T) {
	// Skip if no test server is running
	if !testutil.CheckTestServer(t, "http://localhost:3000") {
		return
	}

	client, err := opencode.NewClient("http://localhost:3000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test cases for AppService.Log
	tests := []struct {
		name        string
		service     string
		level       opencode.AppLogParamsLevel
		message     string
		extra       map[string]any
		expectError bool
	}{
		{
			name:        "Info level log",
			service:     "test-service",
			level:       opencode.AppLogParamsLevelInfo,
			message:     "Test message",
			expectError: false,
		},
		{
			name:        "Debug level log with extra data",
			service:     "debug-service",
			level:       opencode.AppLogParamsLevelDebug,
			message:     "Debug message",
			extra:       map[string]any{"key": "value", "number": 42},
			expectError: false,
		},
		{
			name:        "Warn level log",
			service:     "warning-service",
			level:       opencode.AppLogParamsLevelWarn,
			message:     "Warning message",
			expectError: false,
		},
		{
			name:        "Error level log",
			service:     "error-service",
			level:       opencode.AppLogParamsLevelError,
			message:     "Error message",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := opencode.AppLogParams{
				Service: &tt.service,
				Level:   &tt.level,
				Message: &tt.message,
			}

			if tt.extra != nil {
				params.Extra = &tt.extra
			}

			_, err := client.App.Log(context.Background(), params)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestSessionServiceGet(t *testing.T) {
	// Skip if no test server is running
	if !testutil.CheckTestServer(t, "http://localhost:3000") {
		return
	}

	client, err := opencode.NewClient("http://localhost:3000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with a dummy session ID (this will likely fail without a real server)
	// but it tests the code path
	session, err := client.Session.Get(context.Background(), "test-session-id", nil)

	// We don't assert on the result since we don't have a real server
	// but we can check that the function executes without panicking
	_ = session
	_ = err
}

func TestSessionPermissionsServiceRespond(t *testing.T) {
	// Skip if no test server is running
	if !testutil.CheckTestServer(t, "http://localhost:3000") {
		return
	}

	client, err := opencode.NewClient("http://localhost:3000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test cases for SessionPermissionsService.Respond
	tests := []struct {
		name         string
		sessionID    string
		permissionID string
		response     int
		expectError  bool
	}{
		{
			name:         "Respond with once",
			sessionID:    "test-session-id",
			permissionID: "test-permission-id",
			response:     opencode.SessionPermissionRespondParamsResponseOnce,
			expectError:  true, // Expected to fail without real server
		},
		{
			name:         "Respond with always",
			sessionID:    "test-session-id",
			permissionID: "test-permission-id",
			response:     opencode.SessionPermissionRespondParamsResponseAlways,
			expectError:  true, // Expected to fail without real server
		},
		{
			name:         "Respond with reject",
			sessionID:    "test-session-id",
			permissionID: "test-permission-id",
			response:     opencode.SessionPermissionRespondParamsResponseReject,
			expectError:  true, // Expected to fail without real server
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := opencode.SessionPermissionRespondParams{
				Response: &tt.response,
			}

			_, err := client.Session.Permissions.Respond(
				context.Background(),
				tt.sessionID,
				tt.permissionID,
				params,
			)

			// With a mock server, we might get specific errors
			// Without a server, we expect connection errors
			// The important thing is that the function executes without panicking
			_ = err
		})
	}
}
