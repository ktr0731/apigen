package apigen

import "errors"

var (
	// ErrInvalidDefinition represents the error is caused by definition misconfiguration.
	ErrInvalidDefinition = errors.New("invalid definition")
	// ErrUnimplemented represents the provided definition contains unsupported things (such as HTTP method).
	ErrUnimplemented = errors.New("unimplemented")
)
