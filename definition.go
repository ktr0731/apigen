package apigen

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

type Definition struct {
	Methods []*Method
}

func (d *Definition) validate() error {
	for _, m := range d.Methods {
		if !isIdent(m.Service) {
			return fmt.Errorf("service name '%s' should satisfy identifier name spec: %w", m.Service, ErrInvalidDefinition)
		}
		if !isIdent(m.Method) {
			return fmt.Errorf("method name '%s' should satisfy identifier name spec: %w", m.Method, ErrInvalidDefinition)
		}
		if m.Request == nil {
			return fmt.Errorf("Request field should not be nil: %w", ErrInvalidDefinition)
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
