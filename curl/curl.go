package curl

import (
	"context"
	"fmt"
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

		return newRequest(ctx, u, &flags)
	}
}

func newRequest(ctx context.Context, url *url.URL, flags *flags) (*http.Request, error) {
	if flags.request != http.MethodGet && flags.request != http.MethodPost {
		return nil, fmt.Errorf("unsupported method %s: %w", flags.request, apigen.ErrInvalidDefinition)
	}

	req, err := http.NewRequestWithContext(ctx, flags.request, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request, err = '%s': %w", err, apigen.ErrInvalidDefinition)
	}

	for _, val := range flags.headers {
		sp := strings.SplitN(val, ":", 2)
		req.Header.Add(strings.TrimSpace(sp[0]), strings.TrimSpace(sp[1]))
	}

	return req, nil
}
