package apigen

import (
	"io"
	"net/http"
)

// Option represents an option for Generate.
type Option func(*runner)

// WithHTTPClient specifies the HTTP client for invoking HTTP requests to know API structure.
// Default is *http.DefaultClient.
func WithHTTPClient(c *http.Client) Option {
	return func(r *runner) {
		r.client = c
	}
}

// WithWriter specifies the destination writer for generated files. Default is stdout.
func WithWriter(w io.Writer) Option {
	return func(r *runner) {
		r.writer = w
	}
}

// WithPackage specifies the generated file's package name. Default is main.
func WithPackage(name string) Option {
	return func(r *runner) {
		r.pkg = name
	}
}
