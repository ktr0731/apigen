package curl

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/ktr0731/apigen"
	"github.com/mattn/go-shellwords"
	"github.com/morikuni/failure"
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
			return nil, failure.Translate(err, apigen.ErrInvalidUsage, failure.Context{"cmd": cmd})
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
			return nil, failure.Translate(err, apigen.ErrInvalidUsage, failure.Context{"cmd": cmd})
		}

		// "curl" and URL.
		if fs.NArg() > 2 {
			return nil, failure.New(apigen.ErrInvalidUsage, failure.Message("URL must be specified only one"))
		}

		u, err := url.Parse(fs.Arg(1))
		if err != nil {
			return nil, failure.New(apigen.ErrInvalidUsage, failure.Context{"url": fs.Arg(1)})
		}

		return newRequest(ctx, u, &flags)
	}
}

func newRequest(ctx context.Context, url *url.URL, flags *flags) (*http.Request, error) {
	if flags.request != http.MethodGet && flags.request != http.MethodPost {
		return nil, failure.New(apigen.ErrInvalidUsage, failure.Message("unsupported method"))
	}

	req, err := http.NewRequestWithContext(ctx, flags.request, url.String(), nil) // TODO
	if err != nil {
		return nil, failure.Translate(err, apigen.ErrInvalidUsage, failure.Context{"method": flags.request})
	}

	for _, val := range flags.headers {
		sp := strings.SplitN(val, ":", 2)
		req.Header.Add(strings.TrimSpace(sp[0]), strings.TrimSpace(sp[1]))
	}

	return req, nil
}
