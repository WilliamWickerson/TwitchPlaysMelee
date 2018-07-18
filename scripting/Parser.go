package scripting

import (
	"scripting/token"
	"scripting/AST"
	"errors"
	"fmt"
)

type Parser interface {
	Parse() AST.Script;
}

type parser struct {
	scanner Scanner;
	currToken token.Token;
}

func NewParser(s Scanner) Parser {
	return &parser{scanner : s};
}

func (p *parser) Parse() AST.Script {
	p.currToken = p.scanner.NextToken();
	return p.script();
}

func contains(types []token.Type, ty token.Type) bool {
	for _, t := range types {
		if t == ty {
			return true;
		}
	}
	return false;
}

func (p *parser) match(t token.Type) error {
	fmt.Printf("%d\n", p.currToken.Type());
	defer func(){p.currToken = p.scanner.NextToken()}();;
	if p.currToken.Type() == t {
		return nil;
	} else {
		return errors.New("Error, expected token type: " + string(t) + ", but got: " + string(p.currToken.Type()));
	}
}

func (p *parser) matchArray(types []token.Type) error {
	fmt.Printf("%d\n", p.currToken.Type());
	defer func(){p.currToken = p.scanner.NextToken()}();;
	if contains(types, p.currToken.Type()) {
		return nil;
	} else {
		return errors.New("Error, expected token type: " + string(p.currToken.Type()) + ", not in list");
	}
}

var (
	firstCommand = []token.Type{token.KW_PRESS, token.KW_UNPRESS, token.KW_STICK, token.KW_CSTICK, token.IDENTIFIER};
	firstButtonCommand = []token.Type{token.KW_PRESS, token.KW_UNPRESS};
	buttonMinusLR = []token.Type{token.KW_A, token.KW_B, token.KW_X, token.KW_Y, token.KW_Z, token.KW_START,
	                       token.KW_DLEFT, token.KW_DRIGHT, token.KW_DUP, token.KW_DDOWN};
	buttonLR = []token.Type{token.KW_L, token.KW_R};
	button = append(buttonMinusLR, buttonLR...);
	firstStickCommand = []token.Type{token.KW_STICK, token.KW_CSTICK};
	firstMacroCommand = []token.Type{token.IDENTIFIER};
	inputDirection = []token.Type{token.INTLITERAL, token.FLOATLITERAL};
	namedDirections = []token.Type{token.KW_CENTER, token.KW_LEFT, token.KW_RIGHT, token.KW_UP, token.KW_DOWN};
	firstMacroInput = append([]token.Type{token.INTLITERAL, token.FLOATLITERAL}, namedDirections...);
)

func (p *parser) script() AST.Script {
	commands := make([]AST.Command,0);
	commands = append(commands, p.command());
	for p.currToken.Type() == token.SEMICOLON {
		p.match(token.SEMICOLON);
		if contains(firstCommand, p.currToken.Type()) {
			commands = append(commands, p.command());
		}
	}
	return AST.Script{commands};
}

func (p *parser) command() AST.Command {
	if contains(firstButtonCommand, p.currToken.Type()) {
		return p.buttonCommand();
	} else if contains(firstStickCommand, p.currToken.Type()) {
		return p.stickCommand();
	} else if contains(firstMacroCommand, p.currToken.Type()) {
		return p.macroCommand();
	} else {
		return nil;
	}
}

func (p *parser) buttonCommand() AST.Command {
	command := p.currToken.Type();
	p.matchArray(firstButtonCommand);
	but := p.currToken.Type();
	if contains(buttonMinusLR, but) {
		p.matchArray(buttonMinusLR);
		dur := p.duration();
		return AST.NewButtonCommand(command, but, dur);
	} else if contains(buttonLR, but) {
		p.matchArray(buttonLR);
		var val float64;
		if p.currToken.Type() == token.FLOATLITERAL {
			val,_ = p.currToken.Float();
			p.match(token.FLOATLITERAL);
		} else if command == token.KW_PRESS {
			val = 1;
		} else if command == token.KW_UNPRESS {
			val = 0;
		}
		dur := p.duration();
		return AST.NewSliderCommand(command, but, val, dur);
	}
	//How did we get here?
	return nil;
}

func (p *parser) stickCommand() AST.Command {
	stick := p.currToken.Type();
	p.matchArray(firstStickCommand);
	dir := p.direction();
	dur := p.duration();
	return AST.NewStickCommand(stick, dir, dur);
}

func (p *parser) macroCommand() AST.Command {
	p.match(token.IDENTIFIER);
	for contains(firstMacroInput, p.currToken.Type()) {
		//do stuff with input
	}
	return nil;
}

func (p *parser) duration() AST.Duration {
	frames := make([]AST.Frames, 0);
    frames = append(frames, p.frames());
	for p.currToken.Type() == token.COMMA {
		p.match(token.COMMA);
		frames = append(frames, p.frames());
	}
	return AST.NewDuration(frames);
}

func (p *parser) frames() AST.Frames {
	start,_ := p.currToken.Integer();
	p.match(token.INTLITERAL);
	if (p.currToken.Type() == token.HYPHEN) {
		p.match(token.HYPHEN);
		end,_ := p.currToken.Integer();
		p.match(token.INTLITERAL);
		return AST.Frames{start, end};
	} else {
		return AST.Frames{start, start};
	}
}

func (p *parser) direction() AST.Direction {
	tilt := false;
	if p.currToken.Type() == token.KW_TILT {
		p.match(token.KW_TILT)
		tilt = true;
	}
	if p.currToken.Type() == token.OPENPAREN {
		p.match(token.OPENPAREN);
		x,_ := p.currToken.Float();
		p.matchArray(inputDirection);
		p.match(token.COMMA);
		y,_ := p.currToken.Float();
		p.matchArray(inputDirection);
		p.match(token.CLOSEPAREN);
		return AST.NewDirectionPair(x, y, tilt);
	} else if contains(namedDirections, p.currToken.Type()) {
		directions := make([]token.Type, 0);
		directions = append(directions, p.currToken.Type());
		p.matchArray(namedDirections);
		for p.currToken.Type() == token.HYPHEN {
			p.match(token.HYPHEN);
			directions = append(directions, p.currToken.Type());
			p.matchArray(namedDirections);
		}
		return AST.NewDirection(directions, tilt);
	}
	return nil;
}
