package scripting

type Scanner interface {
	NextToken() Token;
}

type scanner struct {
	Scanner;
	tokens []Token;
}

func NewScanner(s string) Scanner {
	return nil;
}

func (s scanner) NextToken() {

}
