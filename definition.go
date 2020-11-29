package apigen

import (
	"context"
	"net/http"
	"strings"
	"unicode"

	"github.com/morikuni/failure"
)

type Definition struct {
	Methods []*Method
}

func (d *Definition) validate() error {
	for _, m := range d.Methods {
		if !isIdent(m.Service) {
			return failure.New(
				ErrInvalidDefinition,
				failure.Messagef("service name '%s' should satisfy identifier name spec", m.Service),
			)
		}
		if !isIdent(m.Method) {
			return failure.New(
				ErrInvalidDefinition,
				failure.Messagef("method name '%s' should satisfy identifier name spec", m.Method),
			)
		}
		if m.Request == nil {
			return failure.New(
				ErrInvalidDefinition,
				failure.Message("Request field should not be nil"),
			)
		}
	}

	return nil
}

type RequestFunc func(context.Context) (*http.Request, error)

type Method struct {
	Service string
	Method  string
	Request RequestFunc
}

func isIdent(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		// https://golang.org/ref/spec#identifier
		return !(unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r))
	}) == -1
}
