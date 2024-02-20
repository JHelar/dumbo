package template

import (
	"log"

	"github.com/JHelar/dumbo/internal/lex"
	"github.com/JHelar/dumbo/internal/parser"
)

type ComponentRenderer func(children string, props map[string]interface{}) string
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

func (dumbo *Dumbo) Parse(templateString string) *DumboTemplate {
	content := []byte(templateString)
	lexer := lex.NewLexer(content)
	parser := parser.NewParser(lexer)

	expr, err := parser.Parse()
	if err != nil {
		log.Panic(err.Error())
	}

	return &DumboTemplate{
		expr:  expr,
		dumbo: dumbo,
	}
}

func (dumbo *Dumbo) AddComponent(componentName string, renderer ComponentRenderer) {
	dumbo.components[componentName] = DumboComponent{
		name:     componentName,
		renderer: renderer,
	}
}
