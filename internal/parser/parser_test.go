package parser_test

import (
	"log"
	"testing"

	"github.com/JHelar/dumbo/internal/lex"
	"github.com/JHelar/dumbo/internal/parser"
)

func TestSimpleElement(t *testing.T) {
	content := []byte("<div>")
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		log.Print(expression.String())
	}
}

func TestElementWithAttributes(t *testing.T) {
	content := []byte(`<div hx-swap="delete" hx-select="body">`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		log.Print(expression.String())
	}
}
