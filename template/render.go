package template

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/JHelar/dumbo/internal/parser"
)

func renderAttributeValue(expression parser.Expression, data any) any {
	switch t := expression.(type) {
	case parser.AttributeReference:
		dataValue := reflect.ValueOf(data)
		switch dataValue.Kind().String() {
		case "map":
			return dataValue.MapIndex(reflect.ValueOf(t.PropName)).Interface()
		case "struct":
			return dataValue.FieldByName(t.PropName).Interface()
		default:
			log.Panicf("Unsupported attribute value type!")
			return nil
		}
	case parser.String:
		return t.Value
	case parser.Number:
		return t.Value
	default:
		log.Panicf("Unsupported attribute expression!")
		return nil
	}
}

func (dumbo *Dumbo) renderExpr(expression parser.Expression, data any) string {
	switch t := expression.(type) {
	case parser.Element:
		attributes := map[string]any{}
		for _, attribute := range t.Attributes {
			attributes[attribute.Name] = renderAttributeValue(attribute.Value, data)
		}

		children := ""
		for _, child := range t.Children {
			children += dumbo.renderExpr(child, data)
		}

		component, ok := dumbo.components[t.TagName]
		if ok {
			return component.renderer(children, attributes)
		}

		htmlAttributes := []string{}
		for key, value := range attributes {
			htmlAttributes = append(htmlAttributes, fmt.Sprintf("%s=\"%s\"", key, value))
		}
		if len(htmlAttributes) > 0 {
			return fmt.Sprintf("<%s %s>%s</%s>", t.TagName, strings.Join(htmlAttributes, " "), children, t.TagName)
		}
		return fmt.Sprintf("<%s>%s</%s>", t.TagName, children, t.TagName)
	case parser.String:
		return t.Value
	case parser.AttributeReference:
		value := renderAttributeValue(t, data)
		switch tValue := value.(type) {
		case int:
			return fmt.Sprintf("%d", tValue)
		case string:
			return tValue
		default:
			return "[Object object]"
		}
	}

	return ""
}

func (template *DumboTemplate) Render(data any) string {
	result := ""
	for _, expression := range template.expr {
		result += template.dumbo.renderExpr(expression, data)
	}
	return result
}
