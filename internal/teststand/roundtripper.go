package teststand

import "net/http"

// RoundTripperFunc is a test helper to stub net/http transports with a function.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
