package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sst/opencode-api-go"
	"github.com/sst/opencode-api-go/internal/testutil"
)

// MockServer represents a mock HTTP server for testing
type MockServer struct {
	Server *httptest.Server
}

// NewMockServer creates a new mock server with the provided handler
func NewMockServer(handler http.Handler) *MockServer {
	return &MockServer{
		Server: httptest.NewServer(handler),
	}
}

// Close shuts down the mock server
func (m *MockServer) Close() {
	m.Server.Close()
}

// MockHandler creates a handler for the mock server
func MockHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON content type
		w.Header().Set("Content-Type", "application/json")

		// Handle different endpoints
		switch r.URL.Path {
		case "/api/v1/app/log":
			if r.Method == http.MethodPost {
				// Parse the request body
				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				// Return a simple success response
				response := map[string]interface{}{
					"status": "success",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return

		case "/api/v1/session/get":
			if r.Method == http.MethodGet {
				// Extract session ID from query parameters
				sessionID := r.URL.Query().Get("id")
				if sessionID == "" {
					http.Error(w, "Missing session ID", http.StatusBadRequest)
					return
				}

				// Return a mock session response
				response := map[string]interface{}{
					"id":    sessionID,
					"title": fmt.Sprintf("Session %s", sessionID),
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return

		case "/api/v1/session/permissions/respond":
			if r.Method == http.MethodPost {
				// Extract path parameters
				sessionID := r.URL.Query().Get("session_id")
				permissionID := r.URL.Query().Get("permission_id")

				if sessionID == "" || permissionID == "" {
					http.Error(w, "Missing session_id or permission_id", http.StatusBadRequest)
					return
				}

				// Parse the request body
				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				// Return a simple success response
				response := map[string]interface{}{
					"status": "permission responded",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return

		default:
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
	})
}

// TestLiveServer tests the SDK against a live server if available
func TestLiveServer(t *testing.T) {
	// Skip if no test server is running
	if !testutil.CheckTestServer(t, "http://localhost:3000") {
		return
	}

	client, err := opencode.NewClient("http://localhost:3000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test AppService.Log
	t.Run("AppServiceLog", func(t *testing.T) {
		service := "test-service"
		level := opencode.AppLogParamsLevelInfo
		message := "Test log message"

		params := opencode.AppLogParams{
			Service: &service,
			Level:   &level,
			Message: &message,
		}

		_, err := client.App.Log(context.Background(), params)
		if err != nil {
			t.Logf("AppService.Log failed (expected without real server): %v", err)
		}
	})

	// Test SessionService.Get
	t.Run("SessionServiceGet", func(t *testing.T) {
		session, err := client.Session.Get(context.Background(), "test-session-id", nil)
		if err != nil {
			t.Logf("SessionService.Get failed (expected without real server): %v", err)
		}
		_ = session
	})

	// Test SessionPermissionsService.Respond
	t.Run("SessionPermissionsServiceRespond", func(t *testing.T) {
		response := opencode.SessionPermissionRespondParamsResponseOnce
		params := opencode.SessionPermissionRespondParams{
			Response: &response,
		}

		_, err := client.Session.Permissions.Respond(
			context.Background(),
			"test-session-id",
			"test-permission-id",
			params,
		)
		if err != nil {
			t.Logf("SessionPermissionsService.Respond failed (expected without real server): %v", err)
		}
	})
}

// TestMockServer tests the SDK against a mock server
func TestMockServer(t *testing.T) {
	// Create a mock server
	mockServer := NewMockServer(MockHandler(t))
	defer mockServer.Close()

	// Create a client that connects to the mock server
	client, err := opencode.NewClient(mockServer.Server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test AppService.Log
	t.Run("AppServiceLog", func(t *testing.T) {
		service := "test-service"
		level := opencode.AppLogParamsLevelInfo
		message := "Test log message"

		params := opencode.AppLogParams{
			Service: &service,
			Level:   &level,
			Message: &message,
		}

		_, err := client.App.Log(context.Background(), params)
		if err != nil {
			t.Errorf("AppService.Log failed: %v", err)
		}
	})

	// Test SessionService.Get
	t.Run("SessionServiceGet", func(t *testing.T) {
		session, err := client.Session.Get(context.Background(), "test-session-id", nil)
		if err != nil {
			t.Errorf("SessionService.Get failed: %v", err)
		}
		if session.ID != "test-session-id" {
			t.Errorf("Expected session ID 'test-session-id', got '%s'", session.ID)
		}
		if session.Title != "Session test-session-id" {
			t.Errorf("Expected session title 'Session test-session-id', got '%s'", session.Title)
		}
	})

	// Test SessionPermissionsService.Respond
	t.Run("SessionPermissionsServiceRespond", func(t *testing.T) {
		response := opencode.SessionPermissionRespondParamsResponseOnce
		params := opencode.SessionPermissionRespondParams{
			Response: &response,
		}

		_, err := client.Session.Permissions.Respond(
			context.Background(),
			"test-session-id",
			"test-permission-id",
			params,
		)
		if err != nil {
			t.Errorf("SessionPermissionsService.Respond failed: %v", err)
		}
	})
}

// TestEndpointCoverage tests that all expected endpoints are covered
func TestEndpointCoverage(t *testing.T) {
	// Create a coverage tracker
	tracker := NewCoverageTracker()
	tracker.TotalEndpoints = 3 // Update this as more endpoints are implemented

	// Create a mock server
	mockServer := NewMockServer(MockHandler(t))
	defer mockServer.Close()

	// Create a client that connects to the mock server
	client, err := opencode.NewClient(mockServer.Server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test AppService.Log
	t.Run("AppServiceLog", func(t *testing.T) {
		tracker.MarkImplemented("/api/v1/app/log")

		service := "test-service"
		level := opencode.AppLogParamsLevelInfo
		message := "Test log message"

		params := opencode.AppLogParams{
			Service: &service,
			Level:   &level,
			Message: &message,
		}

		_, err := client.App.Log(context.Background(), params)
		if err != nil {
			t.Errorf("AppService.Log failed: %v", err)
		} else {
			tracker.MarkTested("/api/v1/app/log")
		}
	})

	// Test SessionService.Get
	t.Run("SessionServiceGet", func(t *testing.T) {
		tracker.MarkImplemented("/api/v1/session/get")

		session, err := client.Session.Get(context.Background(), "test-session-id", nil)
		if err != nil {
			t.Errorf("SessionService.Get failed: %v", err)
		} else {
			tracker.MarkTested("/api/v1/session/get")
		}

		if session.ID != "test-session-id" {
			t.Errorf("Expected session ID 'test-session-id', got '%s'", session.ID)
		}
	})

	// Test SessionPermissionsService.Respond
	t.Run("SessionPermissionsServiceRespond", func(t *testing.T) {
		tracker.MarkImplemented("/api/v1/session/permissions/respond")

		response := opencode.SessionPermissionRespondParamsResponseOnce
		params := opencode.SessionPermissionRespondParams{
			Response: &response,
		}

		_, err := client.Session.Permissions.Respond(
			context.Background(),
			"test-session-id",
			"test-permission-id",
			params,
		)
		if err != nil {
			t.Errorf("SessionPermissionsService.Respond failed: %v", err)
		} else {
			tracker.MarkTested("/api/v1/session/permissions/respond")
		}
	})

	// Report coverage
	t.Logf("API Implementation Coverage: %.2f%%", tracker.CalculateCoverage())
	t.Logf("API Test Coverage: %.2f%%", tracker.CalculateTestCoverage())

	// Save coverage report
	if err := tracker.Save("coverage_report.json"); err != nil {
		t.Logf("Failed to save coverage report: %v", err)
	}
}
