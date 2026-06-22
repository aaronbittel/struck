package main

import (
	"fmt"
	"reflect"
)

type ParseErrorKind int

const (
	ErrMissingValue ParseErrorKind = iota
	ErrInvalidValue
	ErrMissingRequired
)

type ParamKind int

const (
	KindFlag ParamKind = iota
	KindPositional
)

type ParseError struct {
	kind ParseErrorKind

	argStr       string
	expectedType reflect.Kind

	paramKind ParamKind
	name      string

	err error
}

func (p ParseError) Error() string {
	switch p.kind {
	case ErrMissingValue:
		return fmt.Sprintf("missing value for flag %q", p.name)
	case ErrInvalidValue:
		switch p.paramKind {
		case KindFlag:
			return fmt.Sprintf("invalid value %q for flag %q: not a valid %s", p.argStr, p.name, p.expectedType)
		case KindPositional:
			return fmt.Sprintf("invalid value %q for arg %q: not a valid %s", p.argStr, p.name, p.expectedType)
		default:
			return fmt.Sprintf("invalid value %q: not a valid %s", p.argStr, p.expectedType)
		}
	case ErrMissingRequired:
		switch p.paramKind {
		case KindFlag:
			return fmt.Sprintf("required flag %q not provided", p.name)
		case KindPositional:
			return fmt.Sprintf("missing required argument: %s", p.name)
		default:
			return fmt.Sprintf("missing required value")
		}
	default:
		return fmt.Sprintf("%s", p.err.Error())
	}
}
