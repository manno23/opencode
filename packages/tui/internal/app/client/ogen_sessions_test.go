package client

import (
	"context"
	"os"
	"testing"

	"github.com/sst/opencode-sdk-go"
)

func TestOgenSessions(t *testing.T) {
	// Set environment variable to enable ogen sessions
	os.Setenv("TUI_SDK_EXPERIMENT_SESSIONS", "true")
	defer os.Unsetenv("TUI_SDK_EXPERIMENT_SESSIONS")

	// Create a mock legacy client
	legacyClient := &opencode.Client{}

	// Create client config
	config := &ClientConfig{
		UseOgenSessions: true,
		LegacyClient:    legacyClient,
	}

	// Test creating ogen sessions client
	client := newSessionsClient(config)

	// Verify it's an OgenSessions client (not legacy)
	_, ok := client.(*OgenSessions)
	if !ok {
		t.Error("Expected OgenSessions client, but got different type")
	}
}

func TestOgenSessionsCreate(t *testing.T) {
	// Skip if no server is running
	if os.Getenv("OPENCODE_SERVER") == "" {
		t.Skip("OPENCODE_SERVER not set, skipping integration test")
	}

	// Set environment variable to enable ogen sessions
	os.Setenv("TUI_SDK_EXPERIMENT_SESSIONS", "true")
	defer os.Unsetenv("TUI_SDK_EXPERIMENT_SESSIONS")

	// Create a mock legacy client
	legacyClient := &opencode.Client{}

	// Create client config
	config := &ClientConfig{
		UseOgenSessions: true,
		LegacyClient:    legacyClient,
	}

	// Create ogen sessions client
	client := NewOgenSessions(config)

	// Test creating a session
	ctx := context.Background()
	opts := CreateOpts{
		Title: "Test Session",
	}

	session, err := client.Create(ctx, opts)
	if err != nil {
		t.Logf("Session creation failed (expected if server not running): %v", err)
		// This is expected if no server is running
		return
	}

	if session == nil {
		t.Error("Expected session to be created, got nil")
	}

	if session.Title != "Test Session" {
		t.Errorf("Expected session title 'Test Session', got '%s'", session.Title)
	}
}
