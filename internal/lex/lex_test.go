package lex_test

import (
	"testing"

	"github.com/JHelar/dumbo/internal/lex"
)

func TestTokenSpace(t *testing.T) {
	content := []byte("    spaces")
	lexer := lex.NewLexer(content)

	token, _ := lexer.Next()
	if token.Kind != lex.TokenSpace {
		t.Errorf("Expected token space got: %s", token.Kind)
	}

	if len(token.Content) != 4 {
		t.Errorf("Expected correct amount of spaces got: %d", len(token.Content))
	}
}

func TestTokenSymbol(t *testing.T) {
	content := []byte(`
	test
	
	
	comp
	`)
	lexer := lex.NewLexer(content)

	token, _ := lexer.Next()
	if token.Kind != lex.TokenSymbol {
		t.Errorf("Expected token symbol got: %s", token.Kind)
	}

	if token.Content != "test" {
		t.Errorf("Expected correct content got: %s", token.Content)
	}
}

func TestAst(t *testing.T) {
	content := []byte("{comp:children}")
	lexer := lex.NewLexer(content)

	token1, _ := lexer.Next()
	if token1.Kind != lex.TokenOpenCurly {
		t.Errorf("Expected TokenOpenCurly got: %s", token1.Kind)
	}

	token2, _ := lexer.Next()
	if token2.Kind != lex.TokenSymbol {
		t.Errorf("Expected TokenSymbol got: %s", token2.Kind)
	}

	token3, _ := lexer.Next()
	if token3.Kind != lex.TokenCloseCurly {
		t.Errorf("Expected TokenCloseCurly got: %s", token2.Kind)
	}
}
