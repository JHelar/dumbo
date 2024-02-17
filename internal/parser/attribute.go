package parser

import (
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

type Attribute struct {
	name  string
	value Expression
}

func parseAttributes(lexer *lex.Lexer) ([]Expression, *ParserError) {
	attributes := []Expression{}
	for {
		if nextToken, err := lexer.Peak(); nextToken.Kind == lex.TokenSymbol {
			attribute, err := parseAttribute(lexer)
			if err != nil {
				return nil, err
			}

			attributes = append(attributes, attribute)

			err = skipSpace(lexer)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, newError(InternalErr, "")
		} else {
			break
		}
	}

	return attributes, nil
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
