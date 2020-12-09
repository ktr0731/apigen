package curl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ktr0731/apigen"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/pflag"
)

type flags struct {
	headers    []string
	request    string
	data       string
	compressed bool
}

func ParseCommand(cmd string) apigen.RequestFunc {
	return func(ctx context.Context) (*http.Request, error) {
		args, err := shellwords.Parse(cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to parse command '%s', err = '%s': %w", cmd, err, apigen.ErrInvalidDefinition)
		}

		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}

		var flags flags
		fs := pflag.NewFlagSet("curl", pflag.ContinueOnError)
		fs.StringArrayVarP(&flags.headers, "header", "H", nil, "")
		fs.StringVarP(&flags.request, "request", "X", http.MethodGet, "")
		fs.StringVar(&flags.data, "data-binary", "", "")
		fs.StringVarP(&flags.data, "data", "d", "", "")
		fs.BoolVar(&flags.compressed, "compressed", false, "")

		if err := fs.Parse(args); err != nil {
			return nil, fmt.Errorf("failed to parse curl flags, err = '%s': %w", err, apigen.ErrInvalidDefinition)
		}

		// "curl" and URL.
		if fs.NArg() > 2 {
			return nil, fmt.Errorf("URL must be specified only one: %w", apigen.ErrInvalidDefinition)
		}

		u, err := url.Parse(fs.Arg(1))
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL '%s', err = '%s': %w", fs.Arg(1), err, apigen.ErrInvalidDefinition)
		}

		if flags.data != "" && flags.request == http.MethodGet {
			flags.request = http.MethodPost
		}

		return newRequest(ctx, u, &flags)
	}
}

func newRequest(ctx context.Context, url *url.URL, flags *flags) (*http.Request, error) {
	switch flags.request {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
	default:
		return nil, fmt.Errorf("unsupported method %s: %w", flags.request, apigen.ErrInvalidDefinition)
	}

	var body io.Reader
	if flags.data != "" {
		body = strings.NewReader(flags.data)
	}
	req, err := http.NewRequestWithContext(ctx, flags.request, url.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request, err = '%s': %w", err, apigen.ErrInvalidDefinition)
	}

	for _, val := range flags.headers {
		sp := strings.SplitN(val, ":", 2)
		req.Header.Add(strings.TrimSpace(sp[0]), strings.TrimSpace(sp[1]))
	}

	return req, nil
}
