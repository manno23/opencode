package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea/v2"
	flag "github.com/spf13/pflag"
	"github.com/sst/opencode-sdk-go"
	"github.com/sst/opencode-sdk-go/option"
	"github.com/sst/opencode/internal/api"
	"github.com/sst/opencode/internal/app"
	"github.com/sst/opencode/internal/clipboard"
	"github.com/sst/opencode/internal/tui"
	"github.com/sst/opencode/internal/util"
)

var Version = "dev"

func main() {
	version := Version
	if version != "dev" && !strings.HasPrefix(Version, "v") {
		version = "v" + Version
	}

	var model *string = flag.String("model", "", "model to begin with")
	var prompt *string = flag.String("prompt", "", "prompt to begin with")
	var agent *string = flag.String("agent", "", "agent to begin with")
	var sessionID *string = flag.String("session", "", "session ID")
	flag.Parse()

	url := os.Getenv("OPENCODE_SERVER")

	appInfoStr := os.Getenv("OPENCODE_APP_INFO")
	var appInfo opencode.App
	err := json.Unmarshal([]byte(appInfoStr), &appInfo)
	if err != nil {
		slog.Error("Failed to unmarshal app info", "error", err)
		os.Exit(1)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		slog.Error("Failed to stat stdin", "error", err)
		os.Exit(1)
	}

	// Check if there's data piped to stdin
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error("Failed to read stdin", "error", err)
			os.Exit(1)
		}
		stdinContent := strings.TrimSpace(string(stdin))
		if stdinContent != "" {
			if prompt == nil || *prompt == "" {
				prompt = &stdinContent
			} else {
				combined := *prompt + "\n" + stdinContent
				prompt = &combined
			}
		}
	}

	httpClient := opencode.NewClient(option.WithBaseURL(url))

	// Fetch agents from the /agent endpoint
	agentsPtr, err := httpClient.App.Agents(context.Background())
	if err != nil {
		slog.Error("Failed to fetch agents", "error", err)
		os.Exit(1)
	}
	agents := *agentsPtr
	if len(agents) == 0 {
		slog.Error("No agents returned from server")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiHandler := util.NewAPILogHandler(ctx, httpClient, "tui", slog.LevelDebug)
	logger := slog.New(apiHandler)
	slog.SetDefault(logger)

	slog.Debug("TUI launched", "app", appInfoStr, "agents_count", len(agents), "url", url)

	go func() {
		err = clipboard.Init()
		if err != nil {
			slog.Error("Failed to initialize clipboard", "error", err)
		}
	}()

	// Create main context for the application
	app_, err := app.New(ctx, version, appInfo, agents, httpClient, model, prompt, agent, sessionID)
	if err != nil {
		panic(err)
	}

	tuiModel := tui.NewModel(app_).(*tui.Model)
	program := tea.NewProgram(
		tuiModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// TODO: Implement event streaming with EventSubscribe
	// go func() {
	// 	resp, err := httpClient.EventSubscribe(ctx)
	// 	if err != nil {
	// 		slog.Error("Error subscribing to events", "error", err)
	// 		return
	// 	}
	// 	defer resp.Body.Close()
	// 	// Parse SSE stream
	// 	// ... implementation needed
	// }()

	go api.Start(ctx, program, httpClient)

	// Handle signals in a separate goroutine
	go func() {
		sig := <-sigChan
		slog.Info("Received signal, shutting down gracefully", "signal", sig)
		tuiModel.Cleanup()
		program.Quit()
	}()

	// Run the TUI
	result, err := program.Run()
	if err != nil {
		slog.Error("TUI error", "error", err)
	}

	tuiModel.Cleanup()
	slog.Info("TUI exited", "result", result)
}
