package client

import (
	"context"
	"net/http"
)

type (
	Interceptor func(ctx context.Context, req *http.Request, handler Handler) (*http.Response, error)
	Handler     func(ctx context.Context, req *http.Request) (*http.Response, error)
)

func HeaderInterceptor(h http.Header) Interceptor {
	return func(ctx context.Context, req *http.Request, handler Handler) (*http.Response, error) {
		for k, v := range h {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}

		return handler(ctx, req)
	}
}
