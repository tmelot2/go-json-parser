package main

import (
	// "fmt"
	"testing"
)


func runParserWithStr(s string) (map[string]interface{}, error) {
	result, err := ParseJson(s)
	return result, err
}

func TestParserInvalidJson(t *testing.T) {
	// Test empty string
	_, err := runParserWithStr("")
	if err == nil {
		t.Error("Expected parser error for blank JSON, did not error")
	}

	// Test missing start JSON
	_, err = runParserWithStr("}")
	if err == nil {
		t.Error("Expected parser error on missing end JSON close brace, did not error")
	}

	// Test missing start JSON
	_, err = runParserWithStr(`"hello": "world"}`)
	if err == nil {
		t.Error("Expected parser error on missing end JSON close brace, did not error")
	}

	// Test missing end JSON
	_, err = runParserWithStr("{")
	if err == nil {
		t.Error("Expected parser error on missing end JSON close brace, did not error")
	}

	// Test missing field assignment ":"
	_, err = runParserWithStr(`{"hello" "world"}`)
	if err == nil {
		t.Error("Expected parser error on missing field assignment \":\", did not error")
	}

	// Test missing field separator ","
	_, err = runParserWithStr(`{ "a": 1 "b": 2 }`)
	if err == nil {
		t.Error("Expected parser error on missing field separator \",\", did not error")
	}

	// Test missing quote on key
	_, err = runParserWithStr(`{ a: 1 }`)
	if err == nil {
		t.Error("Expected parser error on missing field separator \",\", did not error")
	}
}
