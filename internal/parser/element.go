package parser

import (
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

type Element struct {
	tagName    string
	children   []Expression
	attributes []Expression
}

func parseElement(lexer *lex.Lexer) (Expression, *ParserError) {
	_, err := nextTokenKind(lexer, lex.TokenLessThen)
	if err != nil {
		return nil, err
	}

	elementNameToken, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}

	err = skipSpace(lexer)
	if err != nil {
		return nil, err
	}

	attributes := []Expression{}

	for {
		if nextToken, err := lexer.Peak(); nextToken.Kind != lex.TokenGreaterThen {
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

	_, err = nextTokenKind(lexer, lex.TokenGreaterThen)
	if err != nil {
		return nil, err
	}

	return Expression(Element{
		tagName:    elementNameToken.Content,
		children:   []Expression{},
		attributes: attributes,
	}), nil
}

func (e Element) String() string {
	attributeValues := ""
	for _, expr := range e.attributes {
		attributeValues += fmt.Sprintf("\n%s", expr.String())
	}

	return fmt.Sprintf(`{
		type: Element
		tagName: %s
		attributes: %s
	}`, e.tagName, attributeValues)
}
