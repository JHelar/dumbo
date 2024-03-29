package lex

import (
	"os"
	"slices"
	"strconv"
)

type Lexer struct {
	currentToken Token
	reader       *runeReader
	currentRow   int
	currentCol   int
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
	TokenComma       TokenKind = ","
	TokenNumber      TokenKind = "number"
	TokenSpace       TokenKind = "space"
	TokenNewline     TokenKind = "newline"
	TokenTab         TokenKind = "tab"
	TokenSymbol      TokenKind = "symbol"
)

var ChildrenSymbol = "children"

var SingleTokens []rune = []rune{'(', ')', '{', '}', '<', '>', '\\', '/', '=', ' ', '"', '\n', ',', '\t'}

type Token struct {
	Content string
	Kind    TokenKind
	Row     int
	Col     int
}

func NewLexerFromFile(file *os.File) *Lexer {
	return &Lexer{
		reader:       newRuneReaderFromFile(file),
		currentToken: Token{},
		currentRow:   0,
		currentCol:   0,
	}
}

func NewLexer(content []byte) *Lexer {
	return &Lexer{
		reader:       newRuneReader(content),
		currentToken: Token{},
		currentRow:   0,
		currentCol:   0,
	}
}

func (lex *Lexer) newToken(kind TokenKind, content string) Token {
	return Token{
		Kind:    kind,
		Content: content,
		Col:     lex.currentCol,
		Row:     lex.currentRow,
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
		case '(', ')', '{', '}', '<', '>', '\\', '/', '=', '"', ',':
			token := lex.newToken(kind, content)
			lex.currentCol++
			return token, nil
		case '\n':
			token := lex.newToken(TokenNewline, content)
			lex.currentRow++
			lex.currentCol = 0
			return token, nil
		case ' ':
			content += lex.takeUntil(func(r rune) bool {
				return r == ' '
			})
			token := lex.newToken(TokenSpace, content)
			lex.currentCol += len(content)
			return token, nil
		case '\t':
			content += lex.takeUntil(func(r rune) bool {
				return r == '\t'
			})
			token := lex.newToken(TokenTab, content)
			lex.currentCol += len(content)
			return token, nil
		default:
			content += lex.takeUntil(func(r rune) bool {
				return !slices.Contains(SingleTokens, r)
			})
			lex.currentCol += len(content)

			if _, numberErr := strconv.Atoi(content); numberErr == nil {
				return lex.newToken(TokenNumber, content), nil
			}

			return lex.newToken(TokenSymbol, content), nil
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
