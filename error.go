package apigen

import "errors"

var (
	ErrInvalidDefinition = errors.New("invalid definition")
	ErrUnimplemented     = errors.New("unimplemented")
)
