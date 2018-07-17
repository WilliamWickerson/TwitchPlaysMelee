package scripting

import (
	"strconv"
)

type Token struct {
	Type Type;
	s string;
}

type Type int;

const (
	INTLITERAL = iota;
	FLOATLITERAL;
	IDENTIFIER;
	COMMA;
	SEMICOLON;
	HYPHEN;
	OPENPAREN;
	CLOSEPAREN;
	KW_PRESS;
	KW_UNPRESS;
	KW_STICK;
	KW_CSTICK;
	KW_CENTER;
	KW_LEFT;
	KW_RIGHT;
	KW_UP;
	KW_DOWN;
	KW_TILT;
	KW_A;
	KW_B;
	KW_X;
	KW_Y;
	KW_Z;
	KW_L;
	KW_R;
	KW_START;
	KW_DLEFT;
	KW_DRIGHT;
	KW_DUP;
	KW_DDOWN;
	EOF;
)

func (t Token) Identifier() string {
	return t.s;
}

func (t Token) Integer() (int, error) {
	return strconv.Atoi(t.s);
}

func (t Token) Float() (float64, error) {
	return strconv.ParseFloat(t.s, 64);
}
