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

type Command struct {
	url *url.URL

	flags *flags
}

func ParseCommand(cmd string) (*Command, error) {
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

	return &Command{url: u, flags: &flags}, nil
}

func (c *Command) Request(ctx context.Context) (*http.Request, error) {
	if c.flags.request != http.MethodGet && c.flags.request != http.MethodPost {
		return nil, failure.New(apigen.ErrInvalidUsage, failure.Message("unsupported method"))
	}

	req, err := http.NewRequestWithContext(ctx, c.flags.request, c.url.String(), nil) // TODO
	if err != nil {
		return nil, failure.Translate(err, apigen.ErrInvalidUsage, failure.Context{"method": c.flags.request})
	}

	for _, val := range c.flags.headers {
		sp := strings.SplitN(val, ":", 2)
		req.Header.Add(strings.TrimSpace(sp[0]), strings.TrimSpace(sp[1]))
	}

	return req, nil
}
