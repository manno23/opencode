package client

import (
	"os"
	"testing"

	"github.com/sst/opencode-sdk-go"
)

func TestSessionsDualImplementation(t *testing.T) {
	// Create a mock legacy client
	legacyClient := &opencode.Client{}

	// Test legacy implementation
	t.Run("Legacy", func(t *testing.T) {
		config := &ClientConfig{
			UseOgenSessions: false,
			LegacyClient:    legacyClient,
		}

		client := newSessionsClient(config)
		_, ok := client.(*LegacySessions)
		if !ok {
			t.Error("Expected LegacySessions client")
		}
	})

	// Test ogen implementation
	t.Run("Ogen", func(t *testing.T) {
		config := &ClientConfig{
			UseOgenSessions: true,
			LegacyClient:    legacyClient,
		}

		client := newSessionsClient(config)
		_, ok := client.(*OgenSessions)
		if !ok {
			t.Error("Expected OgenSessions client")
		}
	})
}

func TestFeatureFlagSwitching(t *testing.T) {
	// Save original env var
	original := os.Getenv("TUI_SDK_EXPERIMENT_SESSIONS")
	defer os.Setenv("TUI_SDK_EXPERIMENT_SESSIONS", original)

	legacyClient := &opencode.Client{}

	// Test with flag disabled
	os.Setenv("TUI_SDK_EXPERIMENT_SESSIONS", "false")
	config := loadConfig(legacyClient)
	if config.UseOgenSessions {
		t.Error("Expected UseOgenSessions to be false")
	}

	// Test with flag enabled
	os.Setenv("TUI_SDK_EXPERIMENT_SESSIONS", "true")
	config = loadConfig(legacyClient)
	if !config.UseOgenSessions {
		t.Error("Expected UseOgenSessions to be true")
	}
}
