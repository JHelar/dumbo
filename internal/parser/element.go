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

func newElementExpression(tagName string, children []Expression, attributes []Expression) Expression {
	return Expression(Element{
		tagName:    tagName,
		children:   children,
		attributes: attributes,
	})
}

func parseElement(lexer *lex.Lexer) (Expression, *ParserError) {
	elementNameToken, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}

	err = skipSpace(lexer)
	if err != nil {
		return nil, err
	}

	tagName := elementNameToken.Content
	attributes, err := parseAttributes(lexer)
	if err != nil {
		return nil, err
	}

	// Check selfclosing element
	closingToken, closingTokenErr := lexer.Next()
	if closingTokenErr != nil {
		return nil, newError(InternalErr, closingTokenErr.Error())
	} else if closingToken.Kind == lex.TokenSlash {
		_, err = nextTokenKind(lexer, lex.TokenGreaterThen)
		if err != nil {
			return nil, err
		}

		return newElementExpression(tagName, []Expression{}, attributes), nil
	}

	if closingToken.Kind != lex.TokenGreaterThen {
		return nil, unexpectedTokenErr(">", closingToken.Content)
	}

	// Parse children
	children := []Expression{}

childLoop:
	for {
		childToken, childTokenErr := lexer.Peak()
		if childTokenErr != nil {
			return nil, newError(InternalErr, childTokenErr.Error())
		}

		switch childToken.Kind {
		case lex.TokenLessThen:
			lexer.Next()
			_, err = peakTokenKind(lexer, lex.TokenSlash)
			if err == nil {
				break childLoop
			}
			// Nested element
			child, err := parseElement(lexer)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		case lex.TokenOpenCurly:
			lexer.Next()
			child, err := parseExpression(lexer)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		case lex.TokenSymbol, lex.TokenNewline:
			child, err := parseString(lexer)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		default:
			return nil, newError(SyntaxErr, fmt.Sprintf("Invalid syntax %s", childToken.Content))
		}
	}

	_, err = nextTokenKind(lexer, lex.TokenSlash)
	if err != nil {
		return nil, err
	}

	endElementName, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}
	if endElementName.Content != tagName {
		return nil, newError(SyntaxErr, fmt.Sprintf("Invalid element end name, expected '%s' got '%s'", tagName, endElementName.Content))
	}

	_, err = nextTokenKind(lexer, lex.TokenGreaterThen)
	if err != nil {
		return nil, err
	}

	return Expression(Element{
		tagName:    tagName,
		children:   children,
		attributes: attributes,
	}), nil
}

func (e Element) String() string {
	attributeValues := ""
	for _, expr := range e.attributes {
		attributeValues += fmt.Sprintf("\n\t%s", expr.String())
	}

	childValues := ""
	for _, expr := range e.children {
		childValues += fmt.Sprintf("\n\t%s", expr.String())
	}

	return fmt.Sprintf(`{
		type: Element
		tagName: %s
		attributes: %s
		children: %s
	}`, e.tagName, attributeValues, childValues)
}
