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
		return t.Kind == lex.TokenSpace || t.Kind == lex.TokenNewline || t.Kind == lex.TokenTab
	})
	if err != nil {
		return nil, internalErr(err)
	}

	_, tokenErr := nextTokenKind(lexer, lex.TokenLessThen)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return parseElement(lexer)
}

func (parser *Parser) Parse() ([]Expression, *ParserError) {
	expressions := []Expression{}

	for {
		if expression, err := parse(parser.lexer); err != nil && err.Err == lex.EOFError {
			return expressions, nil
		} else if err != nil {
			return nil, err
		} else {
			expressions = append(expressions, expression)
		}
	}
}
