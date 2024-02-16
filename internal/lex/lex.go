package lex

import (
	"os"
	"slices"
)

type Lexer struct {
	currentToken Token
	reader       *runeReader
}

type TokenKind string

const (
	TokenOpenParen   TokenKind = "("
	TokenClosedParen TokenKind = ")"
	TokenOpenCurly   TokenKind = "{"
	TokenCloseCurly  TokenKind = "}"
	TokenLessThen    TokenKind = "<"
	TokenGreaterThen TokenKind = ">"
	TokenQuote       TokenKind = "\""
	TokenSlash       TokenKind = "/"
	TokenEquals      TokenKind = "="
	TokenSpace       TokenKind = "space"
	TokenNewline     TokenKind = "newline"
	TokenSymbol      TokenKind = "symbol"
)

var SingleTokens []rune = []rune{'(', ')', '{', '}', '<', '>', '\\', '/', '=', ' ', '"', '\n'}

type Token struct {
	Content string
	Kind    TokenKind
}

func NewLexerFromFile(file *os.File) *Lexer {
	return &Lexer{
		reader:       newRuneReaderFromFile(file),
		currentToken: Token{},
	}
}

func NewLexer(content []byte) *Lexer {
	return &Lexer{
		reader:       newRuneReader(content),
		currentToken: Token{},
	}
}

func (lex *Lexer) takeUntil(pred func(rune) bool) string {
	content := ""

	for {
		if r, err := lex.reader.Peak(); err != nil || !pred(r) {
			break
		}
		content += string(lex.reader.MustNext())
	}

	return content
}

func (lex *Lexer) SkipUntil(pred func(Token) bool) error {
	for {
		val, err := lex.Peak()
		if err != nil {
			return err
		}
		if pred(val) {
			lex.Next()
		} else {
			break
		}
	}

	return nil
}

func (lex *Lexer) Next() (Token, error) {
	if (lex.currentToken != Token{}) {
		token := lex.currentToken
		lex.currentToken = Token{}

		return token, nil
	}

	if r, err := lex.reader.Next(); err != nil {
		return Token{}, err
	} else {
		content := string(r)
		kind := TokenKind(r)

		switch r {
		case '(', ')', '{', '}', '<', '>', '\\', '/', '=', '"':
			return Token{
				Kind:    kind,
				Content: content,
			}, nil
		case ' ':
			content += lex.takeUntil(func(r rune) bool {
				return r == ' '
			})
			return Token{
				Kind:    TokenSpace,
				Content: content,
			}, nil
		default:
			content += lex.takeUntil(func(r rune) bool {
				return !slices.Contains(SingleTokens, r)
			})

			return Token{
				Kind:    TokenSymbol,
				Content: content,
			}, nil
		}
	}
}

func (lex *Lexer) Peak() (Token, error) {
	if (lex.currentToken != Token{}) {
		return lex.currentToken, nil
	}

	token, err := lex.Next()
	if err != nil {
		return Token{}, err
	}

	lex.currentToken = token

	return token, nil
}
