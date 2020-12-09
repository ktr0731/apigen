package apigen

import (
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

// WithOutputDirectory specifies the destination directory for generated files. Default is the current directory.
func WithOutputDirectory(d string) Option {
	return func(r *runner) {
		r.outDir = d
	}
}

func WithPackage(name string) Option {
	return func(r *runner) {
		r.pkg = name
	}
}
