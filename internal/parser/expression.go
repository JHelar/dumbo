package parser

import (
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

func parseExpression(lexer *lex.Lexer) (Expression, *ParserError) {
	var expr Expression
	var err *ParserError

	expressionToken, expressionTokenErr := lexer.Peak()
	if expressionTokenErr != nil {
		return nil, newError(InternalErr, expressionTokenErr.Error())
	}
	switch expressionToken.Kind {
	case lex.TokenLessThen:
		lexer.Next()
		expr, err = parseElement(lexer)
	case lex.TokenQuote:
		lexer.Next()
		expr, err = parseString(lexer)
		if err != nil {
			return nil, err
		}
		_, err = nextTokenKind(lexer, lex.TokenQuote)
	case lex.TokenSymbol:
		expr, err = parseAttributeReference(lexer)
	default:
		return nil, newError(SyntaxErr, fmt.Sprintf("Expression is missing value got: %s", expressionToken.Content))
	}

	if err != nil {
		return nil, err
	}

	_, err = nextTokenKind(lexer, lex.TokenCloseCurly)
	if err != nil {
		return nil, err
	}

	return expr, nil
}
