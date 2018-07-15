package scripting

import (
	"scripting/AST"
)

type Parser interface {
	Parse() AST.Script;
}

type parser struct {
	Parser;
	scanner Scanner;
}

func NewParser(s Scanner) Parser {
	return nil;
}

func (p parser) Parse() AST.Script {
	script := AST.Script{};
	return script;
}
