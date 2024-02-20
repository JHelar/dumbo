package parser

import (
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

type AttributeReference struct {
	PropName string
}

func newAttributeReferenceExpression(value string) Expression {
	return Expression(AttributeReference{
		PropName: value,
	})
}

func (e AttributeReference) String() string {
	return fmt.Sprintf(`{
		type: AttributeReference
		value: %s
	}`, e.PropName)
}

func parseAttributeReference(lexer *lex.Lexer) (Expression, *ParserError) {
	valueToken, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}

	if _, err := peakTokenKind(lexer, lex.TokenCloseCurly); err != nil {
		return nil, err
	}

	return newAttributeReferenceExpression(valueToken.Content), nil
}
