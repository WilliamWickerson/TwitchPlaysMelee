package scripting

import (
	"fmt"
	"strings"
)

type Scanner interface {
	NextToken() Token;
}

type scanner struct {
	Scanner;
	tokens []Token;
	nextIndex int;
}

func NewScanner(s string) Scanner {
	scanner := &scanner{
		nextIndex : -1,
	}
	scanner.scanTokens(s);
	return scanner;
}

func (s *scanner) NextToken() Token {
	s.nextIndex++;
	if (s.nextIndex >= len(s.tokens)) {
		return Token{
			Type : EOF,
		}
	} else {
		return s.tokens[s.nextIndex];
	}
}

//Create an enumeration for the states to make it more clear
type state int;

const (
	none = iota;
	identifier;
	zero;
	integer;
	float;
)

var (
    //Create a map for the keywords to prevent a giant switch statement
	keywordMap = map[string]Type {
		"press" : KW_PRESS,
		"unpress" : KW_UNPRESS,
		"stick" : KW_STICK,
		"cstick" : KW_CSTICK,
		"center" : KW_CENTER,
		"left" : KW_LEFT,
		"right" : KW_RIGHT,
		"up" : KW_UP,
		"down" : KW_DOWN,
		"tilt" : KW_TILT,
		"a" : KW_A,
		"b" : KW_B,
		"x" : KW_X,
		"y" : KW_Y,
		"z" : KW_Z,
		"l" : KW_L,
		"r" : KW_R,
		"start" : KW_START,
		"dleft" : KW_DLEFT,
		"dright" : KW_DRIGHT,
		"dup" : KW_DUP,
		"ddown" : KW_DDOWN,
	}
)

func (sc *scanner) scanTokens(s string) {
	//Convert to lowercase and add EOF character to make parsing simpler by ensuring NONE state by the end
	s = strings.ToLower(s + string(byte(32)));
	sc.tokens = make([]Token, 0);
	//Keep track of current state and start of current token
	var currState state = none;
	var currStart int;
	pos := 0;
	for pos < len(s) {
		ch := s[pos];
		if (currState == none) {
			//We're possibly starting a new token since we're in the none category
			currStart = pos;
			//Letters start an identifier
			if (ch >= 'a' && ch <= 'z') {
				currState = identifier;
			//Numbers start an integer
			} else if (ch >= '1' && ch <= '9') {
				currState = integer;
			}
			//Otherwise switch the options
			switch (ch) {
				case '0':
					currState = zero;
					break;
				case '.':
					currState = float;
					break;
				case '(':
					sc.tokens = append(sc.tokens, Token{OPENPAREN, s[currStart:pos+1]});
					break;
				case ')':
					sc.tokens = append(sc.tokens, Token{CLOSEPAREN, s[currStart:pos+1]});
					break;
				case '-':
					sc.tokens = append(sc.tokens, Token{HYPHEN, s[currStart:pos+1]});
					break;
				case ',':
					sc.tokens = append(sc.tokens, Token{COMMA, s[currStart:pos+1]});
					break;
				case ';':
					sc.tokens = append(sc.tokens, Token{SEMICOLON, s[currStart:pos+1]});
					break;
			}
		} else if (currState == identifier) {
			if !((ch >= 'a' && ch <= 'z') || ch == '-') {
				//Check the keywordMap to see if the identifier is instead a keyword
				keyType, exists := keywordMap[s[currStart:pos]];
				if exists {
					sc.tokens = append(sc.tokens, Token{keyType, s[currStart:pos]});
				//If it's not then use default IDENTIFIER otherwise use the KW_* type
				} else {
					sc.tokens = append(sc.tokens, Token{IDENTIFIER, s[currStart:pos]});
				}
				currState = none;
				continue;
			}
		} else if (currState == zero) {
			//A token starting with 0 can only be 0 or 0.*
			if ch == '.' {
				currState = float;
			} else {
				sc.tokens = append(sc.tokens, Token{INTLITERAL, s[currStart:pos]});
				currState = none;
				continue;
			}
		} else if (currState == integer) {
			//Upon hitting a period the token becomes a float
			if ch == '.' {
				currState = float;
			//Otherwise if it doesn't become a float and not still int, token's over
			} else if !(ch >= '0' && ch <= '9') {
				sc.tokens = append(sc.tokens, Token{INTLITERAL, s[currStart:pos]});
				currState = none;
				continue;
			}
		} else if (currState == float) {
			//Floats can only continue expanding the decimal so when decimal stops add token
			if !(ch >= '0' && ch <= '9') {
				sc.tokens = append(sc.tokens, Token{FLOATLITERAL, s[currStart:pos]});
				currState = none;
				continue;
			}
		} else {
			fmt.Println("Error, scanner entered undefined state!");
		}
		fmt.Printf("%d %c %d\n", pos, ch, currState);
		pos++;
	}
}
