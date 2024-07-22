package jsonParser

import (
	// "errors"
	// "fmt"
	"strconv"
	"testing"

	"tmelot.jsonparser/internal/assert"
)

func runLexerWithStr(s string) ([]Token, error) {
	lexer := newLexer(s)
	// lexer.Debug = true
	result, err := lexer.lex()
	return result, err
}

func TestLexerEmptyJson(t *testing.T) {
	// Test empty
	result, _ := runLexerWithStr("{}")
	assert.Equal(t, len(result), 2, "Expected to lex 2 tokens")

	// Test empty with empty array
	result, _ = runLexerWithStr("{[]}")
	assert.Equal(t, len(result), 4, "Expected to lex 4 tokens")

	// Test lots of nested empty objects & arrays
	result, _ = runLexerWithStr("{[[[{{[{[{[]}]}]}}]]]}")
	assert.Equal(t, len(result), 22, "Expected to lex 22 tokens")

	// Test super empty
	_, err := runLexerWithStr("")
	assert.Nil(t, err, "Expected to lex 0 tokens, errored instead")
}

func TestLexerInvalidJson(t *testing.T) {
	// Test for error on invalid character
	_, err := runLexerWithStr("{/}")
	assert.NotNil(t, err, "Expected an error, did not error")

	// Test for error on invalid character
	_, err = runLexerWithStr("/")
	assert.NotNil(t, err, "Expected an error, did not error")

	// Test for error on invalid characters (single quote)
	_, err = runLexerWithStr(`{'one': 1.111}`)
	assert.NotNil(t, err, "Expected an error, did not error")

	// Test for error on missing quote on key
	_, err = runLexerWithStr(`{one: 1.111}`)
	assert.NotNil(t, err, "Expected an error, did not error")

	// Test for error for unclosed string
	_, err = runLexerWithStr(`{"one: 1.111}`)
	assert.NotNil(t, err, "Expected an error, did not error")

	// Test for error for number and then invalid string
	_, err = runLexerWithStr(`{"one: 1.111 "abcd"}`)
	assert.NotNil(t, err, "Expected an error, did not error")

	// Test for invalid boolean error
	_, err = runLexerWithStr(`{"one": trueee}`)
	assert.NotNil(t, err, "Expected an error, did not error")

	// Test for invalid boolean error
	_, err = runLexerWithStr(`{"one": ffalse}`)
	assert.NotNil(t, err, "Expected an error, did not error")
}

func TestLexerTokenCounts(t *testing.T) {
	// Test with strings
	result, _ := runLexerWithStr(`{"one": "two"}`)
	assert.Equal(t, len(result), 5, "Expected to lex 5 tokens")

	// Test with numbers & strings
	result, _ = runLexerWithStr(`{"one": 1.111, "two": 2.222}`)
	assert.Equal(t, len(result), 9, "Expected to lex 9 tokens")

	// Test mixed whitespace
	result, _ = runLexerWithStr(`{   	 	  "one": 1}`)
	assert.Equal(t, len(result), 5, "Expected to lex 5 tokens")

	// Test newlines
	result, _ = runLexerWithStr(`{
		  "one":
		 1 		}	 	`)
	assert.Equal(t, len(result), 5, "Expected to lex 5 tokens")

	// Test nested tokens
	result, _ = runLexerWithStr(`{"rankings": [{"name": "Smash Ultimate", "rank": 1}, {"name": "Noita", "rank": 2}]}`)
	assert.Equal(t, len(result), 25, "Expected to lex 25 tokens")

	// Test numbers that run to end of unclosed JSON
	// (It's not valid JSON, but the lexer doesn't care, just checking it parses
	// the number token right)
	result, _ = runLexerWithStr(`{"first": 1.00000000000002`)
	floatVal, _ := strconv.ParseFloat(result[len(result)-1].Value, 64)
	assert.Equal(t, floatVal, 1.00000000000002, "")

	// Test with booleans
	result, _ = runLexerWithStr(`{"one": true, "two": false}`)
	assert.Equal(t, len(result), 9, "Expected to lex 9 tokens")
}
