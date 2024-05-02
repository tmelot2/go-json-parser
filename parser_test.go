package main

import (
	"fmt"
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

func TestParserValidJson(t *testing.T) {
	// Test valid JSON with 1 string
	result, _ := runParserWithStr(`{ "a": "1" }`)
	if result["a"] != "1" {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with multiple strings
	result, _ = runParserWithStr(`{ "a": "1", "b": "2" }`)
	if result["a"] != "1" || result["b"] != "2" {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with one int
	result, _ = runParserWithStr(`{ "a": 1 }`)
	if result["a"] != 1 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with multiple ints
	result, _ = runParserWithStr(`{ "a": 1, "b": 2 }`)
	if result["a"] != 1 || result["b"] != 2 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with one float
	result, _ = runParserWithStr(`{ "a": 1.11 }`)
	if result["a"] != 1.11 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with multiple floats
	result, _ = runParserWithStr(`{ "a": 1.11, "b": 2.22 }`)
	if result["a"] != 1.11 || result["b"] != 2.22 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with nested object
	// TODO: Hopefully this casting can be improved in a future design
	result, _ = runParserWithStr(`{ "a": { "b": { "c": 1 } } }`)
	a := result["a"].(map[string]any)
	b := a["b"].(map[string]any)
	c := b["c"].(int)
	if c != 1 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with array
	result, _ = runParserWithStr(`{
		"best_games": [
			{ "rank": 1, "name": "Noita" },
			{ "rank": 2, "name": "Smash Bros" }
		]
	}`)
	arr := result["best_games"].([]any)
	noita := arr[0].(map[string]any)
	smash := arr[1].(map[string]any)
	if noita["rank"] != 1 || smash["rank"] != 2 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}
}
