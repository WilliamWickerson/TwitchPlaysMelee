package scripting

type Scanner interface {
	NextToken() Token;
}

type scanner struct {
	Scanner;
	tokens []Token;
	nextIndex int;
}

func NewScanner(s string) Scanner {
	scanner := &Scanner{
		nextIndex : 0,
	}
	scanner.scanTokens(s);
	return scanner;
}

func (s scanner) NextToken() Token {
	if (nextIndex >= len(tokens)) {
		return Token{
			Type : EOF,
		}
	} else {
		return s.tokens[nextIndex++];
	}
}

func (s *scanner) scanTokens(s string) {
	//TODO: scan the tokens
}
