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

func (parser *Parser) Next() (Expression, *ParserError) {
	token, err := parser.lexer.Peak()
	if err != nil {
		return nil, newError(err, "Error reading next token")
	}

	switch token.Kind {
	case lex.TokenLessThen:
		return parseElement(parser.lexer)
	default:
		return nil, nil
	}
}
