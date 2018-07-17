package scripting

import (
	"testing"
)

func TestKeywords(t *testing.T) {
	sc := NewScanner(`press unpress stick cstick center left right up down tilt
					  A B X Y Z L R START DLEFT DRIGHT DUP DDOWN`);
	//Check that for each token, the type matches the associated type in the keywordMap
	//In other words, all keywords are mapped appropriately
	if token := sc.NextToken(); token.Type != keywordMap[token.Identifier()] {
		t.Errorf("Error, expected token type: %d, but got: %d", keywordMap[token.Identifier()], token.Type);
	}
}

func TestIdentifier(t *testing.T) {
	sc := NewScanner("these;are.some-identifiers");
	if token := sc.NextToken(); token.Type != IDENTIFIER {
		t.Errorf("Error, expected token type: %d, but got: %d", IDENTIFIER, token.Type);
	} else if token.Identifier() != "these" {
		t.Errorf("Error, expected identifier string: %s, but got: %s", "these", token.Identifier());
	}
	sc.NextToken();
	if token := sc.NextToken(); token.Type != IDENTIFIER {
		t.Errorf("Error, expected token type: %d, but got: %d", IDENTIFIER, token.Type);
	} else if token.Identifier() != "are" {
		t.Errorf("Error, expected identifier string: %s, but got: %s", "are", token.Identifier());
	}
	sc.NextToken();
	if token := sc.NextToken(); token.Type != IDENTIFIER {
		t.Errorf("Error, expected token type: %d, but got: %d", IDENTIFIER, token.Type);
	} else if token.Identifier() != "some-identifiers" {
		t.Errorf("Error, expected identifier string: %s, but got: %s", "some-identifiers", token.Identifier());
	}
}

func TestCharacters(t *testing.T) {
	sc := NewScanner(";-(),");
	if token := sc.NextToken(); token.Type != SEMICOLON {
		t.Errorf("Error, expected token type: %d, but got: %d", SEMICOLON, token.Type);
	}
	if token := sc.NextToken(); token.Type != HYPHEN {
		t.Errorf("Error, expected token type: %d, but got: %d", HYPHEN, token.Type);
	}
	if token := sc.NextToken(); token.Type != OPENPAREN {
		t.Errorf("Error, expected token type: %d, but got: %d", OPENPAREN, token.Type);
	}
	if token := sc.NextToken(); token.Type != CLOSEPAREN {
		t.Errorf("Error, expected token type: %d, but got: %d", CLOSEPAREN, token.Type);
	}
	if token := sc.NextToken(); token.Type != COMMA {
		t.Errorf("Error, expected token type: %d, but got: %d", COMMA, token.Type);
	}
}

func TestIntegers(t *testing.T) {
	sc := NewScanner("123 0 0123456789");
	if token := sc.NextToken(); token.Type != INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", INTLITERAL, token.Type);
	} else if val,err := token.Integer(); err != nil || val != 123 {
		t.Errorf("Error, expected int value: %d, but got: %d", 123, val);
	}
	if token := sc.NextToken(); token.Type != INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", INTLITERAL, token.Type);
	} else if val,err := token.Integer(); err != nil || val != 0 {
		t.Errorf("Error, expected int value: %d, but got: %d", 0, val);
	}
	if token := sc.NextToken(); token.Type != INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", INTLITERAL, token.Type);
	} else if val,err := token.Integer(); err != nil || val != 0 {
		t.Errorf("Error, expected int value: %d, but got: %d", 0, val);
	}
	if token := sc.NextToken(); token.Type != INTLITERAL {
		t.Errorf("Error, expected token type: %d, but got: %d", INTLITERAL, token.Type);
	} else if val,err := token.Integer(); err != nil || val != 123456789 {
		t.Errorf("Error, expected int value: %d, but got: %d", 123456789, val);
	}
}
