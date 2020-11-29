package apigen

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

// Definition defines API metadata for code generation.
type Definition struct {
	// Methods defines API methods.
	// Each method must be unique in combination of "Service" and "Method".
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
			return fmt.Errorf("field Request should not be nil: %w", ErrInvalidDefinition)
		}
	}

	return nil
}

// RequestFunc defines a function which instantiates a new *http.Request.
type RequestFunc func(context.Context) (*http.Request, error)

type Method struct {
	// Service defines the name of service which provides APIs through an API server.
	Service string
	// Method defines the name of method which represents an API.
	Method string
	// Request instantiates a new *http.Request. See examples for details.
	Request RequestFunc
}

func isIdent(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		// https://golang.org/ref/spec#identifier
		return !(unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r))
	}) == -1
}
