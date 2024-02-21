package parser

import (
	"fmt"
	"slices"

	"github.com/JHelar/dumbo/internal/lex"
)

type String struct {
	Value string
}

var StringTokens []lex.TokenKind = []lex.TokenKind{
	lex.TokenSymbol,
	lex.TokenSpace,
	lex.TokenNewline,
	lex.TokenComma,
	lex.TokenTab,
	lex.TokenEquals,
	lex.TokenOpenParen,
	lex.TokenClosedParen,
	lex.TokenSlash,
}

func parseString(lexer *lex.Lexer) (Expression, *ParserError) {
	value := ""

	err := nextTokenUntil(lexer, func(t lex.Token) bool {
		if slices.Contains(StringTokens, t.Kind) {
			value += t.Content
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	return Expression(String{
		Value: value,
	}), nil
}

func (e String) String() string {
	return fmt.Sprintf("{\ntype: String\nvalue: %s\n}", e.Value)
}
