package template

import (
	"log"
	"os"

	"github.com/JHelar/dumbo/internal/lex"
)

func ParseFile(name string) {
	file, _ := os.Open(name)
	lexer := lex.NewLexerFromFile(file)

	token, _ := lexer.Next()
	log.Print(token)
}

func Parse(content []byte) {
	lexer := lex.NewLexer(content)

	token, _ := lexer.Next()
	log.Printf("Got kind: %s, %s", token.Kind, token.Content)
}
