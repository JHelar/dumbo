package parser

import (
	"github.com/JHelar/dumbo/internal/lex"
)

type Expression interface {
	String() string
}

type Parser struct {
	lexer *lex.Lexer
}

func NewParser(lexer *lex.Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func parse(lexer *lex.Lexer) (Expression, *ParserError) {
	err := lexer.SkipUntil(func(t lex.Token) bool {
		return t.Kind == lex.TokenSpace || t.Kind == lex.TokenNewline
	})
	if err != nil {
		return nil, newError(InternalErr, err.Error())
	}

	token, err := lexer.Peak()
	if err != nil {
		return nil, newError(err, "Error reading next token")
	}

	switch token.Kind {
	case lex.TokenLessThen:
		lexer.Next()
		return parseElement(lexer)
	case lex.TokenSymbol:
		if isComponentToken(token) {
			lexer.Next()
			return parseComponent(lexer)
		}
	}

	return nil, unexpectedTokenErr("element or component", token.Content)
}

func (parser *Parser) Next() (Expression, *ParserError) {
	return parse(parser.lexer)
}
