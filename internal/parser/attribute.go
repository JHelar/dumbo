package parser

import (
	"fmt"

	"github.com/JHelar/dumbo/internal/lex"
)

type Attribute struct {
	Name  string
	Value Expression
}

func parseAttributes(lexer *lex.Lexer) ([]Attribute, *ParserError) {
	attributes := []Attribute{}
	for {
		if nextToken, err := lexer.Peak(); nextToken.Kind == lex.TokenSymbol {
			attribute, err := parseAttribute(lexer)
			if err != nil {
				return nil, err
			}

			attributes = append(attributes, attribute)

			err = skipSpace(lexer)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, newError(err, "")
		} else {
			break
		}
	}

	return attributes, nil
}

func parseStringValue(lexer *lex.Lexer) (Expression, *ParserError) {
	_, err := nextTokenKind(lexer, lex.TokenQuote)
	if err != nil {
		return Attribute{}, err
	}

	value, err := parseString(lexer)
	if err != nil {
		return Attribute{}, err
	}

	_, err = nextTokenKind(lexer, lex.TokenQuote)
	if err != nil {
		return Attribute{}, err
	}

	return value, nil
}

func parseAttributeReferenceValue(lexer *lex.Lexer) (Expression, *ParserError) {
	var value Expression

	_, err := nextTokenKind(lexer, lex.TokenOpenCurly)
	if err != nil {
		return Attribute{}, err
	}

	if _, err := peakTokenKind(lexer, lex.TokenNumber); err == nil {
		value, err = parseNumber(lexer)
		if err != nil {
			return Attribute{}, err
		}
	} else {
		value, err = parseAttributeReference(lexer)
		if err != nil {
			return Attribute{}, err
		}
	}

	_, err = nextTokenKind(lexer, lex.TokenCloseCurly)
	if err != nil {
		return Attribute{}, err
	}

	return value, nil
}

func parseAttribute(lexer *lex.Lexer) (Attribute, *ParserError) {
	var err *ParserError
	nameToken, err := nextTokenKind(lexer, lex.TokenSymbol)
	if err != nil {
		return Attribute{}, err
	}

	_, err = nextTokenKind(lexer, lex.TokenEquals)
	if err != nil {
		return Attribute{}, err
	}

	var value Expression
	valueTokenType, valueTokenTypeErr := lexer.Peak()
	if valueTokenTypeErr != nil {
		return Attribute{}, newError(ErrInternal, valueTokenTypeErr.Error())
	}

	switch valueTokenType.Kind {
	case lex.TokenQuote:
		value, err = parseStringValue(lexer)
		if err != nil {
			return Attribute{}, err
		}
	case lex.TokenOpenCurly:
		value, err = parseAttributeReferenceValue(lexer)
		if err != nil {
			return Attribute{}, err
		}
	default:
		return Attribute{}, newError(ErrSyntax, fmt.Sprintf("invalid attribute type '%s'", valueTokenType.Kind))
	}

	return Attribute{
		Name:  nameToken.Content,
		Value: value,
	}, nil
}

func (e Attribute) String() string {
	return fmt.Sprintf(`{
		type: Attribute
		name: %s
		value: %s
	}`, e.Name, e.Value.String())
}
