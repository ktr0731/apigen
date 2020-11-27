package apigen

import "github.com/achiku/varfmt"

func public(s string) string { return varfmt.PublicVarName(s) }
