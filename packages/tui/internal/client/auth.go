package client

import (
	"net/http"
)

// authRoundTripper adds authentication headers to HTTP requests
type authRoundTripper struct {
	token string
	base  http.RoundTripper
}

// RoundTrip adds the Authorization header and delegates to the base round tripper
func (art *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if art.token != "" {
		req.Header.Set("Authorization", "Bearer "+art.token)
	}
	return art.base.RoundTrip(req)
}
