package parser_test

import (
	"testing"

	"github.com/JHelar/dumbo/internal/lex"
	"github.com/JHelar/dumbo/internal/parser"
)

func TestSelfClosingElement(t *testing.T) {
	content := []byte("<img />")
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	_, err := parser.Parse()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	}
}

func TestElementWithAttributes(t *testing.T) {
	content := []byte(`<div hx-swap="delete" hx-select="body">Test</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	_, err := parser.Parse()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	}
}

func TestElementWithNumberAttributes(t *testing.T) {
	content := []byte(`<div data-number={9000}>Test</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	_, err := parser.Parse()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	}
}

func TestElementWithNestedElement(t *testing.T) {
	content := []byte(`<div><h1>Heading</h1></div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	_, err := parser.Parse()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	}
}

func TestElementWithExpression(t *testing.T) {
	content := []byte(`<div><h1>{heading}</h1><h2>{subheading}</h2></div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	_, err := parser.Parse()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	}
}

func TestElementWithNewlineChildren(t *testing.T) {
	content := []byte(`<div>
		<p>Paragraph</p>
		<small>On antother line</small>
	</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	_, err := parser.Parse()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	}
}

func TestStringExpression(t *testing.T) {
	content := []byte(`<div>{"Test"}</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	_, err := parser.Parse()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	}
}
