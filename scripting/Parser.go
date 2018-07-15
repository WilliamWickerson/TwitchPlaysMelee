package scripting

type Parser interface {
	Parse() AST.Script;
}

type parser struct {
	Parser;
	scanner Scanner;
}

func NewParser(s Scanner) AST.Script {
	return nil;
}

func (p parser) Parse() Script {
	return nil;
}
