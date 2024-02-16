package parser

import "github.com/JHelar/dumbo/internal/lex"

func nextTokenKind(lexer *lex.Lexer, expectedKind lex.TokenKind) (lex.Token, *ParserError) {
	token, err := lexer.Next()
	if err != nil {
		return lex.Token{}, newError(InternalErr, "")
	}
	if token.Kind != expectedKind {
		return lex.Token{}, unexpectedTokenErr(string(expectedKind), token.Content)
	}

	return token, nil
}

func nextTokenUntil(lexer *lex.Lexer, pred func(lex.Token) bool) *ParserError {
	for {
		token, err := lexer.Peak()
		if err != nil {
			return newError(InternalErr, "")
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
		return t.Kind == lex.TokenSpace
	})
	if err != nil {
		return newError(InternalErr, err.Error())
	}

	return nil
}