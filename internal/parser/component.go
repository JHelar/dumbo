package parser

import (
	"fmt"
	"strings"

	"github.com/JHelar/dumbo/internal/lex"
)

const ComponentSymbol = "comp"
const ComponentChildSymbol = "comp:children"

func isComponentToken(token lex.Token) bool {
	if token.Kind != lex.TokenSymbol {
		return false
	}

	return token.Content == ComponentSymbol
}

func isComponentChildToken(token lex.Token) bool {
	if token.Kind != lex.TokenSymbol {
		return false
	}

	return token.Content == ComponentChildSymbol
}

type Component struct {
	name      string
	propNames []string
	body      Expression
}

func (e Component) String() string {
	propNamesString := strings.Join(e.propNames, ", ")

	return fmt.Sprintf(`{
		type: Component
		name: %s
		propNames: %s
		body: %s
	}`, e.name, propNamesString, e.body.String())
}

/*
comp Name() =
*/
func parseComponent(lexer *lex.Lexer) (Expression, *ParserError) {
	err := skipSpace(lexer)
	if err != nil {
		return nil, err
	}

	componentNameToken, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}

	_, err = nextTokenKind(lexer, lex.TokenOpenParen)
	if err != nil {
		return nil, err
	}

	propNames := []string{}

propNameLoop:
	for {
		err = skipSpace(lexer)
		if err != nil {
			return nil, err
		}

		propNameToken, propNameTokenErr := lexer.Peak()
		if propNameTokenErr != nil {
			return nil, newError(InternalErr, propNameTokenErr.Error())
		}

		switch propNameToken.Kind {
		case lex.TokenSymbol:
			propNames = append(propNames, propNameToken.Content)
			lexer.Next()
		case lex.TokenComma:
			lexer.Next()
		case lex.TokenClosedParen:
			break propNameLoop
		default:
			return nil, unexpectedTokenErr("Symbol or comma", propNameToken.Content)
		}
	}

	lexer.Next()
	err = skipSpace(lexer)
	if err != nil {
		return nil, err
	}

	_, err = nextTokenKind(lexer, lex.TokenEquals)
	if err != nil {
		return nil, err
	}

	err = skipSpace(lexer)
	if err != nil {
		return nil, err
	}

	// Expect element expression
	_, err = nextTokenKind(lexer, lex.TokenLessThen)
	if err != nil {
		return nil, err
	}

	body, err := parseElement(lexer)
	if err != nil {
		return nil, err
	}

	return Expression(Component{
		name:      componentNameToken.Content,
		propNames: propNames,
		body:      body,
	}), nil
}
