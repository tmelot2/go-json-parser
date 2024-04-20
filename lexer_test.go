package main

import "testing"


func runLexerWithStr(s string) ([]Token, error) {
	lexer := newLexer(s)
	result, err := lexer.lex()
	return result, err
}

func TestLexerEmptyJson(t *testing.T) {
	// Test empty
	result, _ := runLexerWithStr("{}")
	if len(result) != 2 {
		t.Errorf("Expected to lex 2 tokens, lexed %d instead", len(result))
	}

	// Test empty with empty array
	result, _ = runLexerWithStr("{[]}")
	if len(result) != 4 {
		t.Errorf("Expected to lex 4 tokens, lexed %d instead", len(result))
	}

	// Test super empty
	_, err := runLexerWithStr("")
	if err != nil {
		t.Error("Expected to lex 0 tokens, errored instead")
	}
}

func TestLexerInvalidJson(t *testing.T) {
	// Test error on invalid
	_, err := runLexerWithStr("{/}")
	if err == nil {
		t.Error("Expected an error, did not error")
	}

	// Test error on invalid
	_, err = runLexerWithStr("/")
	if err == nil {
		t.Error("Expected an error, did not error")
	}
}
