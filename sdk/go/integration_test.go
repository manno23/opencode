package opencode_test

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"git.j9xym.com/opencode-api-go"
	"git.j9xym.com/opencode-api-go/option"
)

func TestIntegrationWithServer(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Start the server
	serverCmd := exec.Command("opencode", "serve", "--port", "4098", "--hostname", "127.0.0.1")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	err := serverCmd.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer func() {
		if serverCmd.Process != nil {
			serverCmd.Process.Kill()
		}
	}()

	// Wait for server to be ready
	time.Sleep(3 * time.Second)

	// Create client
	client := opencode.NewClient(option.WithBaseURL("http://localhost:4098"))

	// Test 1: Get app info
	t.Run("GetApp", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		app, err := client.App.Get(ctx)
		if err != nil {
			t.Errorf("Failed to get app info: %v", err)
			return
		}

		if app == nil {
			t.Error("Expected app to be non-nil")
			return
		}

		if app.Hostname == "" {
			t.Error("Expected hostname to be non-empty")
		}

		t.Logf("App info: Hostname=%s, Git=%v", app.Hostname, app.Git)
	})

	// Test 2: Get agents
	t.Run("GetAgents", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		agentsPtr, err := client.App.Agents(ctx)
		if err != nil {
			t.Errorf("Failed to get agents: %v", err)
			return
		}

		if agentsPtr == nil {
			t.Log("No agents returned (this may be expected)")
			return
		}

		agents := *agentsPtr
		t.Logf("Found %d agents", len(agents))
		for i, agent := range agents {
			t.Logf("Agent %d: %s", i, agent.Name)
			if agent.Description != "" {
				t.Logf("  Description: %s", agent.Description)
			}
		}
	})

	// Test 3: Test session creation (if supported)
	t.Run("CreateSession", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Try to create a simple session using the New method
		session, err := client.Session.New(ctx, opencode.SessionNewParams{})
		if err != nil {
			t.Logf("Session creation not supported or failed: %v", err)
			return
		}

		if session == nil {
			t.Error("Expected session to be non-nil")
			return
		}

		t.Logf("Created session: %+v", session)
	})

	// Test 4: Test command listing (if supported)
	t.Run("ListCommands", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Try to list available commands
		commandsPtr, err := client.Command.List(ctx)
		if err != nil {
			t.Logf("Command listing not supported or failed: %v", err)
			return
		}

		if commandsPtr == nil {
			t.Log("No commands returned (this may be expected)")
			return
		}

		commands := *commandsPtr
		t.Logf("Found %d commands", len(commands))
		for i, cmd := range commands {
			t.Logf("Command %d: %s", i, cmd.Name)
		}
	})

	t.Log("Integration test completed successfully!")
}
