package parser

import "github.com/JHelar/dumbo/internal/lex"

func nextTokenKind(lexer *lex.Lexer, expectedKind lex.TokenKind) (lex.Token, *ParserError) {
	token, err := lexer.Next()
	if err != nil {
		return lex.Token{}, newError(ErrInternal, "")
	}
	if token.Kind != expectedKind {
		return lex.Token{}, unexpectedTokenErr(expectedKind, token)
	}

	return token, nil
}

func peakTokenKind(lexer *lex.Lexer, expectedKind lex.TokenKind) (lex.Token, *ParserError) {
	token, err := lexer.Peak()
	if err != nil {
		return lex.Token{}, newError(ErrInternal, "")
	}
	if token.Kind != expectedKind {
		return lex.Token{}, unexpectedTokenErr(expectedKind, token)
	}

	return token, nil
}

func nextTokenUntil(lexer *lex.Lexer, pred func(lex.Token) bool) *ParserError {
	for {
		token, err := lexer.Peak()
		if err != nil {
			return newError(ErrInternal, "")
		}
		if !pred(token) {
			break
		} else {
			lexer.Next()
		}
	}

	return nil
}

func skipSpace(lexer *lex.Lexer) *ParserError {
	err := lexer.SkipUntil(func(t lex.Token) bool {
		return t.Kind == lex.TokenSpace || t.Kind == lex.TokenTab || t.Kind == lex.TokenNewline
	})
	if err != nil {
		return newError(ErrInternal, err.Error())
	}

	return nil
}
