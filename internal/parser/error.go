package parser

import (
	"errors"
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

var ErrUnexpectedToken = errors.New("unexpected token")
var ErrInternal = errors.New("internal error")
var ErrSyntax = errors.New("syntax error")

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

func unexpectedTokenErr(expected lex.TokenKind, got lex.Token) *ParserError {
	return newError(ErrUnexpectedToken, fmt.Sprintf("[%d:%d] expected: %s got: '%s'(%s)", got.Row, got.Col, expected, got.Content, got.Kind))
}

func internalErr(err error) *ParserError {
	return newError(err, ErrInternal.Error())
}
