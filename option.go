package apigen

import "net/http"

type Option func(*runner)

func WithHTTPClient(c *http.Client) Option {
	return func(r *runner) {
		r.client = c
	}
}
