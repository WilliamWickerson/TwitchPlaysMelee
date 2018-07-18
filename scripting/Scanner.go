package scripting

import (
	"scripting/token"
	"fmt"
	"strings"
)

type Scanner interface {
	NextToken() token.Token;
}

type scanner struct {
	tokens []token.Token;
	nextIndex int;
}

func NewScanner(s string) Scanner {
	scanner := &scanner{
		nextIndex : -1,
	}
	scanner.scanTokens(s);
	return scanner;
}

func (s *scanner) NextToken() token.Token {
	s.nextIndex++;
	if (s.nextIndex >= len(s.tokens)) {
		return token.New(token.EOF, "");
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
	keywordMap = map[string]token.Type {
		"press" : token.KW_PRESS,
		"unpress" : token.KW_UNPRESS,
		"stick" : token.KW_STICK,
		"cstick" : token.KW_CSTICK,
		"center" : token.KW_CENTER,
		"left" : token.KW_LEFT,
		"right" : token.KW_RIGHT,
		"up" : token.KW_UP,
		"down" : token.KW_DOWN,
		"tilt" : token.KW_TILT,
		"a" : token.KW_A,
		"b" : token.KW_B,
		"x" : token.KW_X,
		"y" : token.KW_Y,
		"z" : token.KW_Z,
		"l" : token.KW_L,
		"r" : token.KW_R,
		"start" : token.KW_START,
		"dleft" : token.KW_DLEFT,
		"dright" : token.KW_DRIGHT,
		"dup" : token.KW_DUP,
		"ddown" : token.KW_DDOWN,
	}
)

func (sc *scanner) scanTokens(s string) {
	//Convert to lowercase and add EOF character to make parsing simpler by ensuring NONE state by the end
	s = strings.ToLower(s + string(byte(32)));
	sc.tokens = make([]token.Token, 0);
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
					sc.tokens = append(sc.tokens, token.New(token.OPENPAREN, s[currStart:pos+1]));
					break;
				case ')':
					sc.tokens = append(sc.tokens, token.New(token.CLOSEPAREN, s[currStart:pos+1]));
					break;
				case '-':
					sc.tokens = append(sc.tokens, token.New(token.HYPHEN, s[currStart:pos+1]));
					break;
				case ',':
					sc.tokens = append(sc.tokens, token.New(token.COMMA, s[currStart:pos+1]));
					break;
				case ';':
					sc.tokens = append(sc.tokens, token.New(token.SEMICOLON, s[currStart:pos+1]));
					break;
			}
		} else if (currState == identifier) {
			if !((ch >= 'a' && ch <= 'z') || ch == '-') {
				//Check the keywordMap to see if the identifier is instead a keyword
				keyType, exists := keywordMap[s[currStart:pos]];
				if exists {
					sc.tokens = append(sc.tokens, token.New(keyType, s[currStart:pos]));
				//If it's not then use default IDENTIFIER otherwise use the KW_* type
				} else {
					sc.tokens = append(sc.tokens, token.New(token.IDENTIFIER, s[currStart:pos]));
				}
				currState = none;
				continue;
			}
		} else if (currState == zero) {
			//A token starting with 0 can only be 0 or 0.*
			if ch == '.' {
				currState = float;
			} else {
				sc.tokens = append(sc.tokens, token.New(token.INTLITERAL, s[currStart:pos]));
				currState = none;
				continue;
			}
		} else if (currState == integer) {
			//Upon hitting a period the token becomes a float
			if ch == '.' {
				currState = float;
			//Otherwise if it doesn't become a float and not still int, token's over
			} else if !(ch >= '0' && ch <= '9') {
				sc.tokens = append(sc.tokens, token.New(token.INTLITERAL, s[currStart:pos]));
				currState = none;
				continue;
			}
		} else if (currState == float) {
			//Floats can only continue expanding the decimal so when decimal stops add token
			if !(ch >= '0' && ch <= '9') {
				sc.tokens = append(sc.tokens, token.New(token.FLOATLITERAL, s[currStart:pos]));
				currState = none;
				continue;
			}
		} else {
			fmt.Println("Error, scanner entered undefined state!");
		}
		pos++;
	}
}
