package dumbo

import (
	"log"
	"os"

	"github.com/JHelar/dumbo/internal/lex"
	"github.com/JHelar/dumbo/internal/parser"
)

type ComponentChildren string
type ComponentProps map[string]interface{}
type ComponentRenderer func(children ComponentChildren, props ComponentProps) (DumboElement, error)
type DumboComponent struct {
	name     string
	renderer ComponentRenderer
}

type DumboTemplate struct {
	expr  []parser.Expression
	dumbo *Dumbo
}

type Dumbo struct {
	components map[string]DumboComponent
}

func New() *Dumbo {
	return &Dumbo{
		components: map[string]DumboComponent{},
	}
}

func (dumbo *Dumbo) parse(lexer *lex.Lexer) (*DumboTemplate, error) {
	parser := parser.NewParser(lexer)

	expr, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return &DumboTemplate{
		expr:  expr,
		dumbo: dumbo,
	}, nil
}

func (dumbo *Dumbo) Parse(templateString string) (*DumboTemplate, error) {
	content := []byte(templateString)
	lexer := lex.NewLexer(content)

	return dumbo.parse(lexer)
}

func (dumbo *Dumbo) ParseFile(filename string) (*DumboTemplate, error) {
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		return nil, fileErr
	}
	defer file.Close()

	lexer := lex.NewLexerFromFile(file)

	return dumbo.parse(lexer)
}

func Must(template *DumboTemplate, err error) *DumboTemplate {
	if err != nil {
		log.Panicf("Dumbo error %s", err.Error())
	}
	return template
}

func (dumbo *Dumbo) AddComponent(componentName string, renderer ComponentRenderer) {
	dumbo.components[componentName] = DumboComponent{
		name:     componentName,
		renderer: renderer,
	}
}
