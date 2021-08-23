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
	// Services defines API services and its methods.
	// Each method name must be unique.
	Services map[string][]*Method
}

func (d *Definition) validate() error {
	for service, methods := range d.Services {
		for _, m := range methods {
			if !isIdent(service) {
				return fmt.Errorf("service name '%s' should satisfy identifier name spec: %w", service, ErrInvalidDefinition)
			}
			if !isIdent(m.Name) {
				return fmt.Errorf("method name '%s' should satisfy identifier name spec: %w", m.Name, ErrInvalidDefinition)
			}
			if m.Request == nil {
				return fmt.Errorf("field Request should not be nil: %w", ErrInvalidDefinition)
			}
		}
	}

	return nil
}

// RequestFunc defines a function which instantiates a new *http.Request.
// RequestFunc may return errors wrapping a pre-defined error in apigen.
type RequestFunc func(context.Context) (*http.Request, error)

// Method defines an API method.
type Method struct {
	// Name defines the name of method which represents an API.
	Name string
	// Request instantiates a new *http.Request. See examples for details.
	Request RequestFunc
	// ParamHint specifies path parameters.
	// These will be organized as the request fields. Each parameter must be surrounded by "{}".
	// For example, "/posts/{postID}" is given as a ParamHint, apigen generates the following request type:
	//
	//  type Request struct {
	//    PostID string
	//  }
	//
	// If ParamHint differs the actual path, it will be ignored and never generate any request fields.
	ParamHint string
}

func isIdent(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		// https://golang.org/ref/spec#identifier
		return !(unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r))
	}) == -1
}
