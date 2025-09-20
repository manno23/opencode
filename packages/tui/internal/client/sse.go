package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	api "github.com/sst/opencode-api-go/api"
	ssestream "github.com/sst/opencode-api-go/packages/ssestream"
)

// SSEConfig holds configuration for the SSE client
type SSEConfig struct {
	BaseURL              string
	AuthToken            string
	ReconnectDelay       time.Duration
	MaxReconnectAttempts int
	Timeout              time.Duration
}

// SSEClient manages Server-Sent Events connection
type SSEClient struct {
	config     SSEConfig
	client     *http.Client
	stream     *ssestream.Stream[api.EventListResponse]
	streamMux  sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	eventChan  chan api.EventListResponse
	errorChan  chan error
	connected  bool
	connectMux sync.RWMutex
}

// NewSSEClient creates a new SSE client
func NewSSEClient(config SSEConfig) *SSEClient {
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 5 * time.Second
	}
	if config.MaxReconnectAttempts == 0 {
		config.MaxReconnectAttempts = 10
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	// Add authentication if token provided
	if config.AuthToken != "" {
		httpClient.Transport = &authRoundTripper{
			token: config.AuthToken,
			base:  http.DefaultTransport,
		}
	}

	return &SSEClient{
		config:    config,
		client:    httpClient,
		eventChan: make(chan api.EventListResponse, 100), // Buffered channel
		errorChan: make(chan error, 10),
	}
}

// Connect establishes the SSE connection
func (s *SSEClient) Connect(ctx context.Context) error {
	s.streamMux.Lock()
	defer s.streamMux.Unlock()

	if s.connected {
		return fmt.Errorf("already connected")
	}

	s.ctx, s.cancel = context.WithCancel(ctx)

	// Create HTTP request for event subscription
	req, err := http.NewRequestWithContext(s.ctx, "GET", s.config.BaseURL+"/api/v1/event/subscribe", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to SSE endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("SSE endpoint returned status %d", resp.StatusCode)
	}

	// Create SSE decoder and stream
	decoder := ssestream.NewDecoder(resp)
	if decoder == nil {
		resp.Body.Close()
		return fmt.Errorf("failed to create SSE decoder")
	}

	s.stream = ssestream.NewStream[api.EventListResponse](decoder, nil)
	s.setConnected(true)

	slog.Info("SSE connection established", "url", s.config.BaseURL)

	// Start event processing goroutine
	go s.processEvents()

	return nil
}

// Disconnect closes the SSE connection
func (s *SSEClient) Disconnect() error {
	s.streamMux.Lock()
	defer s.streamMux.Unlock()

	s.setConnected(false)

	if s.cancel != nil {
		s.cancel()
	}

	if s.stream != nil {
		return s.stream.Close()
	}

	close(s.eventChan)
	close(s.errorChan)

	slog.Info("SSE connection disconnected")
	return nil
}

// Events returns a channel for receiving events
func (s *SSEClient) Events() <-chan api.EventListResponse {
	return s.eventChan
}

// Errors returns a channel for receiving errors
func (s *SSEClient) Errors() <-chan error {
	return s.errorChan
}

// IsConnected returns true if the client is connected
func (s *SSEClient) IsConnected() bool {
	s.connectMux.RLock()
	defer s.connectMux.RUnlock()
	return s.connected
}

// setConnected sets the connection status
func (s *SSEClient) setConnected(connected bool) {
	s.connectMux.Lock()
	defer s.connectMux.Unlock()
	s.connected = connected
}

// processEvents processes incoming SSE events
func (s *SSEClient) processEvents() {
	defer func() {
		s.setConnected(false)
	}()

	for s.stream.Next() {
		select {
		case <-s.ctx.Done():
			return
		default:
			event := s.stream.Current()
			select {
			case s.eventChan <- event:
			case <-s.ctx.Done():
				return
			default:
				slog.Warn("Event channel full, dropping event")
			}
		}
	}

	if err := s.stream.Err(); err != nil {
		select {
		case s.errorChan <- fmt.Errorf("SSE stream error: %w", err):
		default:
			slog.Error("SSE stream error", "error", err)
		}
	}
}

// Reconnect attempts to reconnect with exponential backoff
func (s *SSEClient) Reconnect(ctx context.Context) error {
	if s.IsConnected() {
		return nil
	}

	var lastErr error
	delay := s.config.ReconnectDelay

	for attempt := 1; attempt <= s.config.MaxReconnectAttempts; attempt++ {
		slog.Info("Attempting SSE reconnection", "attempt", attempt, "max_attempts", s.config.MaxReconnectAttempts)

		if err := s.Connect(ctx); err != nil {
			lastErr = err
			slog.Warn("SSE reconnection failed", "attempt", attempt, "error", err)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				// Exponential backoff with max delay of 5 minutes
				delay = delay * 2
				if delay > 5*time.Minute {
					delay = 5 * time.Minute
				}
				continue
			}
		}

		slog.Info("SSE reconnection successful", "attempt", attempt)
		return nil
	}

	return fmt.Errorf("failed to reconnect after %d attempts: %w", s.config.MaxReconnectAttempts, lastErr)
}
