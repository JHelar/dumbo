package parser

import (
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

type Attribute struct {
	name  string
	value Expression
}

func parseAttribute(lexer *lex.Lexer) (Expression, *ParserError) {
	nameToken, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}

	_, err = nextTokenKind(lexer, lex.TokenEquals)
	if err != nil {
		return nil, err
	}

	_, err = nextTokenKind(lexer, lex.TokenQuote)
	if err != nil {
		return nil, err
	}

	value, err := parseString(lexer)
	if err != nil {
		return nil, err
	}

	_, err = nextTokenKind(lexer, lex.TokenQuote)
	if err != nil {
		return nil, err
	}

	return Expression(Attribute{
		name:  nameToken.Content,
		value: value,
	}), nil
}

func (e Attribute) String() string {
	return fmt.Sprintf(`{
		type: Attribute
		name: %s
		value: %s
	}`, e.name, e.value.String())
}
