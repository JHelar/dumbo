package dumbo

import (
	"fmt"
	"log"
	"reflect"

	"github.com/JHelar/dumbo/internal/lex"
	"github.com/JHelar/dumbo/internal/parser"
)

type DumboElement struct {
	content string
}

func newElement(content string) DumboElement {
	return DumboElement{
		content: content,
	}
}

func renderAttributeValue(expression parser.Expression, data any) (any, error) {
	switch t := expression.(type) {
	case parser.AttributeReference:
		dataValue := reflect.ValueOf(data)
		var fieldValue reflect.Value
		switch dataValue.Kind().String() {
		case "map":
			fieldValue = dataValue.MapIndex(reflect.ValueOf(t.PropName))
		case "struct":
			fieldValue = dataValue.FieldByName(t.PropName)
		default:
			return nil, fmt.Errorf("unsupported data type: %s", dataValue.Kind().String())
		}

		if (fieldValue == reflect.Value{}) {
			if t.PropName == lex.ChildrenSymbol {
				log.Print("No children field found in data")
				return "", nil
			}
			return nil, fmt.Errorf("field %s was not found in data", t.PropName)
		}
		return fieldValue.Interface(), nil
	case parser.String:
		return t.Value, nil
	case parser.Number:
		return t.Value, nil
	default:
		return nil, fmt.Errorf("unsupported attribute expression: %s", expression.String())
	}
}

func getStringRepresentation(value any) (string, error) {
	switch tValue := value.(type) {
	case int:
		return fmt.Sprintf("%d", tValue), nil
	case string:
		return tValue, nil
	default:
		return reflect.ValueOf(value).String(), nil
	}
}

func (dumbo *Dumbo) renderExpr(expression parser.Expression, data any) (DumboElement, error) {
	switch t := expression.(type) {
	case parser.Element:
		attributes := map[string]any{}
		for _, attribute := range t.Attributes {
			attributeValue, attributeValueErr := renderAttributeValue(attribute.Value, data)
			if attributeValueErr != nil {
				return DumboElement{}, attributeValueErr
			}
			attributes[attribute.Name] = attributeValue
		}

		children := ""
		for _, child := range t.Children {
			childHtml, renderChildErr := dumbo.renderExpr(child, data)
			if renderChildErr != nil {
				return DumboElement{}, renderChildErr
			}
			children += childHtml.content
		}

		component, ok := dumbo.components[t.TagName]
		if ok {
			return component.renderer(ComponentChildren(children), attributes)
		}

		htmlAttributes := ""
		for key, value := range attributes {
			valueStr, valueStrErr := getStringRepresentation(value)
			if valueStrErr != nil {
				return DumboElement{}, valueStrErr
			}

			htmlAttributes += fmt.Sprintf(" %s=\"%s\"", key, valueStr)
		}
		if parser.IsSelfClosingElement(t.TagName) {
			return newElement(fmt.Sprintf("<%s%s>", t.TagName, htmlAttributes)), nil
		}
		return newElement(fmt.Sprintf("<%s%s>%s</%s>", t.TagName, htmlAttributes, children, t.TagName)), nil
	case parser.String:
		return newElement(t.Value), nil
	case parser.AttributeReference:
		value, valueErr := renderAttributeValue(t, data)
		if valueErr != nil {
			return DumboElement{}, valueErr
		}
		stringValue, valueErr := getStringRepresentation(value)
		if valueErr != nil {
			return DumboElement{}, valueErr
		}
		return newElement(stringValue), nil

	default:
		return DumboElement{}, nil
	}
}

func (template *DumboTemplate) Render(data any) (DumboElement, error) {
	content, err := template.renderToString(data)

	if err != nil {
		return DumboElement{}, err
	}

	return newElement(content), nil
}

func (template *DumboTemplate) renderToString(data any) (string, error) {
	result := ""
	for _, expression := range template.expr {
		renderResult, renderErr := template.dumbo.renderExpr(expression, data)
		if renderErr != nil {
			return "", renderErr
		}
		result += renderResult.content
	}
	return result, nil
}
