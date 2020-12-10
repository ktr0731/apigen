package client

import (
	"context"
	"fmt"
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

type Error struct{ StatusCode int }

func (e *Error) Error() string {
	return fmt.Sprintf("code %d", e.StatusCode)
}

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
