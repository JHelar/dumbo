package parser

import (
	"errors"
	"fmt"
)

var UnexpectedTokenErr = errors.New("Unexpected token")
var InternalErr = errors.New("Internal error")
var SyntaxErr = errors.New("Syntax error")

type ParserError struct {
	Err    error
	Reason string
}

func (err *ParserError) Error() string {
	return fmt.Sprintf("%s %s", err.Err.Error(), err.Reason)
}

func newError(err error, reason string) *ParserError {
	return &ParserError{
		Err:    err,
		Reason: reason,
	}
}

func unexpectedTokenErr(expected string, got string) *ParserError {
	return newError(UnexpectedTokenErr, fmt.Sprintf("Expected: '%s' got: '%s'", expected, got))
}
