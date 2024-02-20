package parser

import (
	"fmt"
	"strconv"

	"github.com/JHelar/dumbo/internal/lex"
)

type Number struct {
	Value int
}

func parseNumber(lexer *lex.Lexer) (Expression, *ParserError) {
	value, err := nextTokenKind(lexer, lex.TokenNumber)
	if err != nil {
		return nil, err
	}

	number, numberErr := strconv.Atoi(value.Content)
	if numberErr != nil {
		return nil, newError(ErrInternal, fmt.Sprintf("cannot parse value '%s' to number", value.Content))
	}

	return Expression(Number{
		Value: number,
	}), nil
}

func (e Number) String() string {
	return fmt.Sprintf("{\ntype: Number\nvalue: %d\n}", e.Value)
}
