package parser

import (
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

type String struct {
	value string
}

func parseString(lexer *lex.Lexer) (Expression, *ParserError) {
	value := ""

	err := nextTokenUntil(lexer, func(t lex.Token) bool {
		if t.Kind == lex.TokenSymbol || t.Kind == lex.TokenSpace || t.Kind == lex.TokenNewline || t.Kind == lex.TokenComma {
			value += t.Content
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	return Expression(String{
		value: value,
	}), nil
}

func (e String) String() string {
	return fmt.Sprintf("{\ntype: String\nvalue: %s\n}", e.value)
}
