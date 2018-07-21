package scripting

import (
	"testing"
	"scripting/token"
)

func TestKeywords(t *testing.T) {
	sc := NewScanner(`press unpress stick cstick center left right up down tilt
					  A B X Y Z L R START DLEFT DRIGHT DUP DDOWN`);
	//Check that for each token, the type matches the associated type in the keywordMap
	//In other words, all keywords are mapped appropriately
	if tok := sc.NextToken(); tok.Type() != keywordMap[tok.Identifier()] {
		t.Errorf("Error, expected token type: %d, but got: %d", keywordMap[tok.Identifier()], tok.Type());
	}
}

/*
Tests for identifiers with various seperators including hyphenated words
*/
func TestIdentifier(t *testing.T) {
	sc := NewScanner("these;are.some-identifiers");
	if tok := sc.NextToken(); tok.Type() != token.IDENTIFIER {
		t.Errorf("Error, expected token type: %d, but got: %d", token.IDENTIFIER, tok.Type());
	} else if tok.Identifier() != "these" {
		t.Errorf("Error, expected identifier string: %s, but got: %s", "these", tok.Identifier());
	}
	sc.NextToken();
	if tok := sc.NextToken(); tok.Type() != token.IDENTIFIER {
		t.Errorf("Error, expected token type: %d, but got: %d", token.IDENTIFIER, tok.Type());
	} else if tok.Identifier() != "are" {
		t.Errorf("Error, expected identifier string: %s, but got: %s", "are", tok.Identifier());
	}
	sc.NextToken();
	if tok := sc.NextToken(); tok.Type() != token.IDENTIFIER {
		t.Errorf("Error, expected token type: %d, but got: %d", token.IDENTIFIER, tok.Type());
	} else if tok.Identifier() != "some-identifiers" {
		t.Errorf("Error, expected identifier string: %s, but got: %s", "some-identifiers", tok.Identifier());
	}
}

/*
Tests the single character tokens
*/
func TestCharacters(t *testing.T) {
	sc := NewScanner(";-(),");
	if tok := sc.NextToken(); tok.Type() != token.SEMICOLON {
		t.Errorf("Error, expected token type: %d, but got: %d", token.SEMICOLON, tok.Type());
	}
	if tok := sc.NextToken(); tok.Type() != token.HYPHEN {
		t.Errorf("Error, expected token type: %d, but got: %d", token.HYPHEN, tok.Type());
	}
	if tok := sc.NextToken(); tok.Type() != token.OPENPAREN {
		t.Errorf("Error, expected token type: %d, but got: %d", token.OPENPAREN, tok.Type());
	}
	if tok := sc.NextToken(); tok.Type() != token.CLOSEPAREN {
		t.Errorf("Error, expected token type: %d, but got: %d", token.CLOSEPAREN, tok.Type());
	}
	if tok := sc.NextToken(); tok.Type() != token.COMMA {
		t.Errorf("Error, expected token type: %d, but got: %d", token.COMMA, tok.Type());
	}
}

/*
Tests integer literal tokens
*/
func TestIntegers(t *testing.T) {
	sc := NewScanner("123 0 0123456789");
	if tok := sc.NextToken(); tok.Type() != token.INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.INTLITERAL, tok.Type());
	} else if val,err := tok.Integer(); err != nil || val != 123 {
		t.Errorf("Error, expected int value: %d, but got: %d", 123, val);
	}
	if tok := sc.NextToken(); tok.Type() != token.INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.INTLITERAL, tok.Type());
	} else if val,err := tok.Integer(); err != nil || val != 0 {
		t.Errorf("Error, expected int value: %d, but got: %d", 0, val);
	}
	if tok := sc.NextToken(); tok.Type()!= token.INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.INTLITERAL, tok.Type());
	} else if val,err := tok.Integer(); err != nil || val != 0 {
		t.Errorf("Error, expected int value: %d, but got: %d", 0, val);
	}
	if tok := sc.NextToken(); tok.Type() != token.INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.INTLITERAL, tok.Type());
	} else if val,err := tok.Integer(); err != nil || val != 123456789 {
		t.Errorf("Error, expected int value: %d, but got: %d", 123456789, val);
	}
}

/*
Tests float literal tokens
*/
func TestFloat(t *testing.T) {
	sc := NewScanner("1..23 45 0.6789");
	if tok := sc.NextToken(); tok.Type() != token.FLOATLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.FLOATLITERAL, tok.Type());
	} else if val,err := tok.Float(); err != nil || val != 1 {
		t.Errorf("Error, expected int value: %f, but got: %f", 1.0, val);
	}
	if tok := sc.NextToken(); tok.Type() != token.FLOATLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.FLOATLITERAL, tok.Type());
	} else if val,err := tok.Float(); err != nil || val != .23 {
		t.Errorf("Error, expected int value: %f, but got: %f", .23, val);
	}
	if tok := sc.NextToken(); tok.Type() != token.INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.INTLITERAL, tok.Type());
	} else if val,err := tok.Integer(); err != nil || val != 45 {
		t.Errorf("Error, expected int value: %d, but got: %d", 45, val);
	}
	if tok := sc.NextToken(); tok.Type() != token.FLOATLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", token.FLOATLITERAL, tok.Type());
	} else if val,err := tok.Float(); err != nil || val != .6789 {
		t.Errorf("Error, expected int value: %f, but got: %f", .6789, val);
	}
}
