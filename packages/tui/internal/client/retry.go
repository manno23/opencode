package client

import (
	"math"
	"net/http"
	"time"
)

// retryRoundTripper implements http.RoundTripper with exponential backoff retry logic
type retryRoundTripper struct {
	maxRetries int
	base       http.RoundTripper
}

// RoundTrip executes the HTTP request with retry logic
func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rt.maxRetries; attempt++ {
		resp, err = rt.base.RoundTrip(req)

		// If successful or client error (4xx), don't retry
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// If this is the last attempt, return the error
		if attempt == rt.maxRetries {
			break
		}

		// Calculate backoff delay: baseDelay * 2^attempt
		baseDelay := time.Second
		delay := time.Duration(float64(baseDelay) * math.Pow(2, float64(attempt)))

		// Cap the delay at 30 seconds
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}

		time.Sleep(delay)
	}

	return resp, err
}
