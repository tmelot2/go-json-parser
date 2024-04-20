package main

import (
	"testing"
)


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

func TestLexerTokenCounts(t *testing.T) {
	// Test with strings
	result, _ := runLexerWithStr(`{"one": "two"}`)
	if len(result) != 5 {
		t.Errorf("Expected to lex 5 tokens, lexed %d instead", len(result))
	}

	// Test with numbers & strings
	result, _ = runLexerWithStr(`{"one": 1.111, "two": 2.222}`)
	if len(result) != 9 {
		t.Errorf("Expected to lex 9 tokens, lexed %d instead", len(result))
	}

	// Test mixed whitespace
	result, _ = runLexerWithStr(`{   	 	  "one": 1}`)
	if len(result) != 5 {
		t.Errorf("Expected to lex 5 tokens, lexed %d instead", len(result))
	}

	// Test newlines
	result, _ = runLexerWithStr(`{
		  "one":
		 1 		}	 	`)
	if len(result) != 5 {
		t.Errorf("Expected to lex 5 tokens, lexed %d instead", len(result))
	}

	// Test nested tokens
	result, _ = runLexerWithStr(`{"rankings": [{"name": "Smash Ultimate", "rank": 1}, {"name": "Noita", "rank": 2}]}`)
	if len(result) != 25 {
		t.Errorf("Expected to lex 25 tokens, lexed %d instead", len(result))
	}
}