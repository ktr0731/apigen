package apigen

import "github.com/morikuni/failure"

var (
	ErrInvalidUsage      failure.Code = failure.StringCode("InvalidUsage")
	ErrInvalidDefinition failure.Code = failure.StringCode("InvalidDefinition")
)
