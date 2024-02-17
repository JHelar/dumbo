package parser_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/JHelar/dumbo/internal/lex"
	"github.com/JHelar/dumbo/internal/parser"
)

func TestSimpleElement(t *testing.T) {
	content := []byte("<img />")
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
	content := []byte(`<div hx-swap="delete" hx-select="body">Test</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		fmt.Println(expression.String())
	}
}

func TestElementWithNestedElement(t *testing.T) {
	content := []byte(`<div><h1>Heading</h1></div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		fmt.Println(expression.String())
	}
}

func TestElementWithExpression(t *testing.T) {
	content := []byte(`<div><h1>{heading}</h1><h2>{subheading}</h2></div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		fmt.Println(expression.String())
	}
}

func TestElementWithNewlineChildren(t *testing.T) {
	content := []byte(`<div>
		<p>Paragraph</p>
		<small>On antother line</small>
	</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		fmt.Println(expression.String())
	}
}

func TestStringExpression(t *testing.T) {
	content := []byte(`<div>{"Test"}</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		fmt.Println(expression.String())
	}
}

func TestSimpleComponentExpression(t *testing.T) {
	content := []byte(`comp Test() = <div>{"Test"}</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		fmt.Println(expression.String())
	}
}

func TestComponentExpression(t *testing.T) {
	content := []byte(`comp Test(heading) = <div>{heading}</div>`)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expression, err := parser.Next()
	if err != nil {
		t.Errorf("Should be able to parse, got error: %s", err.Error())
	} else {
		fmt.Println(expression.String())
	}
}
