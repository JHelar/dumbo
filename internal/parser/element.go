package parser

import (
	"fmt"
	"slices"

	"github.com/JHelar/dumbo/internal/lex"
)

type SelfClosingTag string

const (
	SelfClosingArea   SelfClosingTag = "area"
	SelfClosingBase   SelfClosingTag = "base"
	SelfClosingBr     SelfClosingTag = "br"
	SelfClosingCol    SelfClosingTag = "col"
	SelfClosingEmbed  SelfClosingTag = "embed"
	SelfClosingHr     SelfClosingTag = "hr"
	SelfClosingImg    SelfClosingTag = "img"
	SelfClosingInput  SelfClosingTag = "input"
	SelfClosingLink   SelfClosingTag = "link"
	SelfClosingMeta   SelfClosingTag = "meta"
	SelfClosingParam  SelfClosingTag = "param"
	SelfClosingSource SelfClosingTag = "source"
	SelfClosingTrack  SelfClosingTag = "track"
	SelfClosingWbr    SelfClosingTag = "wbr"
)

var SelfClosingTags []SelfClosingTag = []SelfClosingTag{
	SelfClosingArea,
	SelfClosingBase,
	SelfClosingBr,
	SelfClosingCol,
	SelfClosingEmbed,
	SelfClosingHr,
	SelfClosingImg,
	SelfClosingInput,
	SelfClosingLink,
	SelfClosingMeta,
	SelfClosingParam,
	SelfClosingSource,
	SelfClosingTrack,
	SelfClosingWbr,
}

func IsSelfClosingElement(tagName string) bool {
	return slices.Contains(SelfClosingTags, SelfClosingTag(tagName))
}

type Element struct {
	TagName    string
	Children   []Expression
	Attributes []Attribute
}

func newElementExpression(tagName string, children []Expression, attributes []Attribute) Expression {
	return Expression(Element{
		TagName:    tagName,
		Children:   children,
		Attributes: attributes,
	})
}

func parseElement(lexer *lex.Lexer) (Expression, *ParserError) {
	elementNameToken, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}

	err = skipSpace(lexer)
	if err != nil {
		return nil, err
	}

	tagName := elementNameToken.Content
	attributes, err := parseAttributes(lexer)
	if err != nil {
		return nil, err
	}

	// Check self closing element
	closingToken, closingTokenErr := lexer.Next()
	if closingTokenErr != nil {
		return nil, newError(ErrInternal, closingTokenErr.Error())
	}
	// If self closed with trailing slash
	if closingToken.Kind == lex.TokenSlash {
		_, err = nextTokenKind(lexer, lex.TokenGreaterThen)
		if err != nil {
			return nil, err
		}

		return newElementExpression(tagName, []Expression{}, attributes), nil
	}

	// If self closed without trailing slash (only applies to known html elements)
	if closingToken.Kind == lex.TokenGreaterThen && IsSelfClosingElement(tagName) {
		return newElementExpression(tagName, []Expression{}, attributes), nil
	}

	// Element start tag should be closed with a trailing >
	if closingToken.Kind != lex.TokenGreaterThen {
		return nil, unexpectedTokenErr(lex.TokenGreaterThen, closingToken)
	}

	// Parse children
	children := []Expression{}

childLoop:
	for {
		childToken, childTokenErr := lexer.Peak()
		if childTokenErr != nil {
			return nil, newError(ErrInternal, childTokenErr.Error())
		}

		switch childToken.Kind {
		case lex.TokenLessThen:
			lexer.Next()
			_, err = peakTokenKind(lexer, lex.TokenSlash)
			if err == nil {
				break childLoop
			}
			// Nested element
			child, err := parseElement(lexer)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		case lex.TokenOpenCurly:
			lexer.Next()
			child, err := parseExpression(lexer)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		case lex.TokenSymbol, lex.TokenNewline:
			child, err := parseString(lexer)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		default:
			return nil, newError(ErrSyntax, fmt.Sprintf("Invalid syntax %s", childToken.Content))
		}
	}

	_, err = nextTokenKind(lexer, lex.TokenSlash)
	if err != nil {
		return nil, err
	}

	endElementName, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return nil, err
	}
	if endElementName.Content != tagName {
		return nil, newError(ErrSyntax, fmt.Sprintf("Invalid element end name, expected '%s' got '%s'", tagName, endElementName.Content))
	}

	_, err = nextTokenKind(lexer, lex.TokenGreaterThen)
	if err != nil {
		return nil, err
	}

	return Expression(Element{
		TagName:    tagName,
		Children:   children,
		Attributes: attributes,
	}), nil
}

func (e Element) String() string {
	attributeValues := ""
	for _, expr := range e.Attributes {
		attributeValues += fmt.Sprintf("\n\t%s", expr.String())
	}

	childValues := ""
	for _, expr := range e.Children {
		childValues += fmt.Sprintf("\n\t%s", expr.String())
	}

	return fmt.Sprintf(`{
		type: Element
		tagName: %s
		attributes: %s
		children: %s
	}`, e.TagName, attributeValues, childValues)
}
