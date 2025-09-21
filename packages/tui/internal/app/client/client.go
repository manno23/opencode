package client

import (
	"os"

	"github.com/sst/opencode-sdk-go"
)

// Client provides a unified API interface that abstracts between legacy and ogen implementations
// It also exposes pass-through SDK services (App, Session, etc.) for existing callers.
type Client struct {
	// Facade domains
	Sessions SessionsClient
	Messages MessagesClient
	Tools    ToolsClient
	Files    FilesClient
	Health   HealthClient

	// SDK pass-through services used by existing codepaths
	Client  *opencode.Client
	App     *opencode.AppService
	Session *opencode.SessionService
	File    *opencode.FileService
	Command *opencode.CommandService
	Agent   *opencode.AgentService
	Find    *opencode.FindService
}

// ClientConfig holds feature flags and client configuration
type ClientConfig struct {
	// Feature flags for per-domain switching
	UseOgenSessions bool
	UseOgenMessages bool
	UseOgenTools    bool
	UseOgenFiles    bool
	UseOgenHealth   bool

	// Legacy SDK client
	LegacyClient *opencode.Client
}

// NewClient creates a new unified client with per-domain implementations based on feature flags
func NewClient(legacyClient *opencode.Client) *Client {
	config := loadConfig(legacyClient)

	return &Client{
		Sessions: newSessionsClient(config),
		Messages: newMessagesClient(config),
		Tools:    newToolsClient(config),
		Files:    newFilesClient(config),
		Health:   newHealthClient(config),
		// Expose SDK services for compatibility
		Client:  config.LegacyClient,
		App:     config.LegacyClient.App,
		Session: config.LegacyClient.Session,
		File:    config.LegacyClient.File,
		Command: config.LegacyClient.Command,
		Agent:   config.LegacyClient.Agent,
		Find:    config.LegacyClient.Find,
	}
}

// loadConfig reads environment variables and constructs client config
func loadConfig(legacyClient *opencode.Client) *ClientConfig {
	return &ClientConfig{
		UseOgenSessions: os.Getenv("TUI_SDK_EXPERIMENT_SESSIONS") == "true",
		UseOgenMessages: os.Getenv("TUI_SDK_EXPERIMENT_MESSAGES") == "true",
		UseOgenTools:    os.Getenv("TUI_SDK_EXPERIMENT_TOOLS") == "true",
		UseOgenFiles:    os.Getenv("TUI_SDK_EXPERIMENT_FILES") == "true",
		UseOgenHealth:   os.Getenv("TUI_SDK_EXPERIMENT_HEALTH") == "true",
		LegacyClient:    legacyClient,
	}
}

// Domain-specific constructors
func newSessionsClient(config *ClientConfig) SessionsClient {
	if config.UseOgenSessions {
		return NewOgenSessions(config)
	}
	return NewLegacySessions(config.LegacyClient)
}

func newMessagesClient(config *ClientConfig) MessagesClient {
	if config.UseOgenMessages {
		// TODO: return NewOgenMessages(config)
		return NewLegacyMessages(config.LegacyClient)
	}
	return NewLegacyMessages(config.LegacyClient)
}

func newToolsClient(config *ClientConfig) ToolsClient {
	if config.UseOgenTools {
		// TODO: return NewOgenTools(config)
		return NewLegacyTools(config.LegacyClient)
	}
	return NewLegacyTools(config.LegacyClient)
}

func newFilesClient(config *ClientConfig) FilesClient {
	if config.UseOgenFiles {
		// TODO: return NewOgenFiles(config)
		return NewLegacyFiles(config.LegacyClient)
	}
	return NewLegacyFiles(config.LegacyClient)
}

func newHealthClient(config *ClientConfig) HealthClient {
	if config.UseOgenHealth {
		// TODO: return NewOgenHealth(config)
		return NewLegacyHealth(config.LegacyClient)
	}
	return NewLegacyHealth(config.LegacyClient)
}
