package scripting

import (
	"scripting/token"
	"scripting/AST"
	"errors"
	"strings"
	"strconv"
)

type Parser interface {
	Parse() []AST.Command;
}

type parser struct {
	scanner Scanner;
	currToken token.Token;
}

func NewParser(s Scanner) Parser {
	return &parser{scanner : s};
}

func (p *parser) Parse() []AST.Command {
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

func (p *parser) consume() {
	p.currToken = p.scanner.NextToken();
}

func (p *parser) match(t token.Type) error {
	if p.currToken.Type() == t {
		p.consume();
		return nil;
	} else {
		return errors.New("Error, expected token type: " + string(t) + ", but got: " + string(p.currToken.Type()));
	}
}

func (p *parser) matchArray(types []token.Type) error {
	if contains(types, p.currToken.Type()) {
		p.consume();
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

/* Script := Command (semicolon (Command | ε))* (semicolon | ε) */
func (p *parser) script() []AST.Command {
	commands := make([]AST.Command,0);
	//Parse first command and add if there's no error
	if command, err := p.command(); err == nil {
		commands = append(commands, command);
	//If there is an error then consume tokens until hitting a semicolon signifying next command
	} else {
		for p.currToken.Type() != token.SEMICOLON && p.currToken.Type() != token.EOF {
			p.consume();
		}
	}
	//while the current token type is semicolon, keep parsing commands
	for p.currToken.Type() == token.SEMICOLON {
		p.consume();
		//Same procedure as first but that wasn't preceeded by semicolon
		if command, err := p.command(); err == nil {
			commands = append(commands, command);
		} else {
			for p.currToken.Type() != token.SEMICOLON && p.currToken.Type() != token.EOF {
				p.consume();
			}
		}
	}
	return commands;
}

/* Command := ButtonCommand | StickCommand | MacroCommand */
func (p *parser) command() (AST.Command, error) {
	//We have an LL(1) grammar so just check the next token to see what to do
	if contains(firstButtonCommand, p.currToken.Type()) {
		return p.buttonCommand();
	} else if contains(firstStickCommand, p.currToken.Type()) {
		return p.stickCommand();
	} else if contains(firstMacroCommand, p.currToken.Type()) {
		return p.macroCommand();
	} else {
		return nil, errors.New("Unrecognized command type");
	}
}

/* ButtonCommand := unpress Button Duration | press ButtonMinusLR Duration |
                    press ButtonLR (float-literal | ε) Duration */
func (p *parser) buttonCommand() (AST.Command, error) {
	//Command should be either press or unpress
	command := p.currToken.Type();
	if err := p.matchArray(firstButtonCommand); err != nil {
		return nil, err;
	}
	but := p.currToken.Type();
	//If the button isn't a slider then we use first type
	if contains(buttonMinusLR, but) {
		p.consume();
		dur, err := p.duration();
		if err != nil {
			return nil, err;
		}
		return AST.NewButtonCommand(command, but, dur), nil;
	//If the button is a slider we need additional logic for the analogue value
	} else if contains(buttonLR, but) {
		p.consume();
		var val float64;
		//Value was specified so use it
		if p.currToken.Type() == token.FLOATLITERAL {
			val,_ = p.currToken.Float();
			p.consume();
		//Otherwise press implies 1.0 analogue press and unpress implies 0.0 analogue press
		} else if command == token.KW_PRESS {
			val = 1;
		} else if command == token.KW_UNPRESS {
			val = 0;
		}
		//Try to parse duration and go from there
		dur, err := p.duration();
		if err != nil {
			return nil, err;
		}
		return AST.NewSliderCommand(command, but, val, dur), nil;
	}
	//How did we get here?
	return nil, errors.New("Unrecognized slider command type");
}

/* StickCommand := (stick | cstick) Direction Duration */
func (p *parser) stickCommand() (AST.Command, error) {
	//Stick should be either 'stick' or 'cstick'
	stick := p.currToken.Type();
	if err := p.matchArray(firstStickCommand); err != nil {
		return nil, err;
	}
	//Parse the direction looking for errors
	dir, err := p.direction();
	if err != nil {
		return nil, err;
	}
	//Parse the duration looking for errors
	dur, err := p.duration();
	if err != nil {
		return nil, err;
	}
	return AST.NewStickCommand(stick, dir, dur), nil;
}

/* MacroCommand := macro-identifier (int-literal | float-literal | namedDirection)* (ε | int-literal) */
func (p *parser) macroCommand() (AST.Command, error) {
	name := p.currToken.Text();
	p.match(token.IDENTIFIER);
	//Get the inputs to the macro
	inputs := make([]token.Token, 0);
	for contains(firstMacroInput, p.currToken.Type()) {
		inputs = append(inputs, p.currToken);
		p.consume();
	}
	//The delay is ambiguous, so we need to check with and without it
	delay := 0;
	text := AST.GetMacro(name, len(inputs));
	if text == "" {
		text = AST.GetMacro(name, len(inputs) - 1);
		if len(inputs) > 0 && inputs[len(inputs) - 1].Type() == token.INTLITERAL && text != "" {
			delay,_ = inputs[len(inputs) - 1].Integer();
			inputs = inputs[:len(inputs) - 1];
		} else {
			return nil, errors.New("Macro does not exist");
		}
	}
	//Replace the #i's in the text with the corresponding input
	for i,token := range inputs {
		text = strings.Replace(text, "#" + strconv.Itoa(i+1), token.Text(), -1);
	}
	//Parse the macro for it's constituent commands and return
	parser := NewParser(NewScanner(text));
	return AST.NewMacroCommand(parser.Parse(), delay), nil;
}

/* Duration := ε | Frames (comma Frames)* */
func (p *parser) duration() (AST.Duration, error) {
	//Make an array to hold the Frames structs
	frames := make([]AST.Frames, 0);
	//Must contain at least one so parse the first
	frame, err := p.frames();
	if err != nil {
		return nil, err;
	}
    frames = append(frames, frame);
    //Commas signify there being more so keeping parsing Frames
	for p.currToken.Type() == token.COMMA {
		p.consume();
		frame, err := p.frames();
		if err != nil {
			return nil, err;
		}
		frames = append(frames, frame);
	}
	return AST.NewDuration(frames), nil;
}

/* Frames := Expression (ε | hyphen Expression) */
func (p *parser) frames() (AST.Frames, error) {
	//Duration can equal ε, but it's simpler to take care of that case here
	if (p.currToken.Type() == token.SEMICOLON || p.currToken.Type() == token.EOF) {
		return AST.Frames{1,1}, nil;
	}
	//Make sure frames starts with an integer
	start, err := p.expression();
	if err != nil {
		return AST.Frames{}, err;
	}
	//If there is a hyphen then it expresses a range so parse that
	if (p.currToken.Type() == token.HYPHEN) {
		p.consume();
		end, err := p.expression();
		if err != nil {
			return AST.Frames{}, err;
		}
		return AST.Frames{start, end}, nil;
	//Otherwise just return simple 1 frame range
	} else {
		return AST.Frames{start, start}, nil;
	}
}

/* Expression := int-literal (ε | plus int-literal) */
func (p *parser) expression() (int, error) {
	val,_ := p.currToken.Integer();
	if err := p.match(token.INTLITERAL); err != nil {
		return -1, err;
	}
	//The only expression is plus a number so look for that
	if p.currToken.Type() == token.PLUS {
		p.consume();
		val2,_ := p.currToken.Integer();
		if err := p.match(token.INTLITERAL); err != nil {
			return -1, err;
		}
		return val + val2, nil;
	}
	return val, nil;
}

/* Direction := (tilt | ε) (PairDirection | NamedDirections) */
func (p *parser) direction() (AST.Direction, error) {
	//Tilt is a modifier at the start, check if it exists
	tilt := false;
	if p.currToken.Type() == token.KW_TILT {
		p.match(token.KW_TILT)
		tilt = true;
	}
	//If next token is first of PairDirection parse that
	if p.currToken.Type() == token.OPENPAREN {
		return p.pairDirection(tilt);
	//Otherwise if it's first of NameDirections parse that isntead
	} else if contains(namedDirections, p.currToken.Type()) {
		return p.namedDirections(tilt);
	}
	return nil, errors.New("Unrecognized direction type");
}

/* PairDirection := (openParen (ε | hyphen) (int-literal | float-literal) comma (ε | hyphen) (int-literal | float-literal) closeParen) */
func (p *parser) pairDirection(tilt bool) (AST.Direction, error) {
	if err := p.match(token.OPENPAREN); err != nil {
		return nil, err;
	}
	//Since inputs can be negative, check whether a hyphen precedes
	firstSign := 1.0;
	if p.currToken.Type() == token.HYPHEN {
		p.consume();
		firstSign = -1.0;
	}
	x,_ := p.currToken.Float();
	x = firstSign * x;
	if err := p.matchArray(inputDirection); err != nil {
		return nil, err;
	}
	if err := p.match(token.COMMA); err != nil {
		return nil, err;
	}
	//Same as previous for the second int/float input
	secondSign := 1.0;
	if p.currToken.Type() == token.HYPHEN {
		p.consume();
		secondSign = -1.0;
	}
	y,_ := p.currToken.Float();
	y = secondSign * y;
	if err := p.matchArray(inputDirection); err != nil {
		return nil, err;
	}
	if err := p.match(token.CLOSEPAREN); err != nil {
		return nil, err;
	}
	return AST.NewDirectionPair(x, y, tilt), nil;
}

/* NamedDirections := NamedDirection NamedDirection* */
func (p *parser) namedDirections(tilt bool) (AST.Direction, error) {
	//Create a slice for the directions
	directions := make([]token.Type, 0);
	//There must be at least one direction so much that
	directions = append(directions, p.currToken.Type());
	if err := p.matchArray(namedDirections); err != nil {
		return nil, err;
	}
	//While there are more directions, keep appending them
	for contains(namedDirections, p.currToken.Type()) {
		directions = append(directions, p.currToken.Type());
		p.consume();
	}
	return AST.NewDirection(directions, tilt), nil;
}
