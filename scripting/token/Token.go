package token

import (
	"strconv"
)

type Token interface {
	Type() Type;
	Text() string;
	Integer() (int, error);
	Float() (float64, error);
}

type token struct {
	t Type;
	s string;
}

func New(t Type, s string) Token {
	return token{t, s};
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
	PLUS;
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

func (t token) Type() Type {
	return t.t;
}

func (t token) Text() string {
	return t.s;
}

func (t token) Integer() (int, error) {
	return strconv.Atoi(t.s);
}

func (t token) Float() (float64, error) {
	return strconv.ParseFloat(t.s, 64);
}
