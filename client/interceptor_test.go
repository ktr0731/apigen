package client_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ktr0731/apigen/client"
)

func TestConvertStatusCodeToErrorInterceptor(t *testing.T) {
	cases := map[string]struct {
		code    int
		err     error
		wantErr bool
	}{
		"1XX":               {code: http.StatusContinue, wantErr: true},
		"2XX":               {code: http.StatusCreated},
		"3XX":               {code: http.StatusFound, wantErr: true},
		"4XX":               {code: http.StatusForbidden, wantErr: true},
		"5XX":               {code: http.StatusGatewayTimeout, wantErr: true},
		"error is returned": {err: errors.New("err"), wantErr: true},
	}

	for name, c := range cases {
		c := c

		t.Run(name, func(t *testing.T) {
			res, err := client.ConvertStatusCodeToErrorInterceptor()(
				context.Background(),
				nil,
				func(context.Context, *http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: c.code}, c.err
				},
			)
			if !c.wantErr {
				if err != nil {
					t.Fatalf("should not return an error, but got '%s'", err)
				}
				if res == nil {
					t.Errorf("res should not be nil")
				}
				return
			}

			if err == nil {
				t.Fatalf("should return an error, but got nil")
			}

			if c.err != nil {
				return
			}

			var e *client.Error
			if !errors.As(err, &e) {
				t.Fatalf("errors.As should not be failed")
			}

			if e.StatusCode != c.code {
				t.Errorf("want code %d, but got %d", e.StatusCode, c.code)
			}
		})
	}
}
