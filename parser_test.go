package main

import (
	"testing"
)


func runParserWithStr(s string) (string, error) {
	lexer := newLexer(s)
	// lexer.Debug = true
	tokens, _ := lexer.lex()
	parser := newParser(tokens)
	// parser.Debug = true
	result, err := parser.Parse()
	return result, err
}

func TestParserInvalidJson(t *testing.T) {
	// Test for error on missing open bracket to start JSON
	_, err := runParserWithStr(`"one": "two"}`)
	if err == nil {
		t.Error("Expected to error on missing open bracket to start JSON, it did not error")
	}

	// Test for error on missing close bracket to end JSON
	_, err = runParserWithStr(`{"one": "two"`)
	if err == nil {
		t.Error("Expected to error on missing close bracket to end JSON, it did not error")
	}
}
