package apigen

import "net/http"

// Option represents an option for Generate.
type Option func(*runner)

// WithHTTPClient specifies the HTTP client for invoking HTTP requests to know API structure.
func WithHTTPClient(c *http.Client) Option {
	return func(r *runner) {
		r.client = c
	}
}
