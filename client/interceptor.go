package client

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// Interceptor intercepts the execution of a method call on the client.
	Interceptor func(ctx context.Context, req *http.Request, handler Handler) (*http.Response, error)
	// Handler represents the actual invoker for handling requests.
	Handler func(ctx context.Context, req *http.Request) (*http.Response, error)
)

// HeaderInterceptor append headers into each request.
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

// Error represents client's error.
type Error struct{ StatusCode int }

// Error returns string representation for the error.
func (e *Error) Error() string {
	return fmt.Sprintf("code %d", e.StatusCode)
}

// ConvertStatusCodeToErrorInterceptor converts non-2XX status code to an *Error.
// If you want more fine-grained handling, define your interceptor instead of using this interceptor.
func ConvertStatusCodeToErrorInterceptor() Interceptor {
	return func(ctx context.Context, req *http.Request, handler Handler) (*http.Response, error) {
		res, err := handler(ctx, req)
		if err != nil {
			return res, err
		}

		if res.StatusCode/100 != 2 {
			return nil, &Error{StatusCode: res.StatusCode}
		}

		return res, err
	}
}
