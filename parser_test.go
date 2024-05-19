package main

import (
	"fmt"
	"testing"
)

func runParserWithStr(s string) (*JsonValue, error) {
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

	// Test invalid object trailing comma
	_, err = runParserWithStr(`{ "a": 1, "b": 2, }`)
	if err == nil {
		t.Error("Expected parser error on missing field separator \",\", did not error")
	}

	// Test invalid array trailing comma
	_, err = runParserWithStr(`{ "a": [1,2,] }`)
	if err == nil {
		t.Error("Expected parser error on trailing array comma")
	}

	// Test invalid array trailing commaa
	_, err = runParserWithStr(`{ "a": [1,2,,] }`)
	if err == nil {
		t.Error("Expected parser error on multiple trailing array commas")
	}
}

func TestParserValidJson(t *testing.T) {
	// Test valid JSON with 1 string
	result, _ := runParserWithStr(`{ "a": "1" }`)
	val, _ := result.GetString("a")
	if val != "1" {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with multiple strings
	result, _ = runParserWithStr(`{ "a": "1", "b": "2" }`)
	resultA, _ := result.GetString("a")
	resultB, _ := result.GetString("b")
	if resultA != "1" || resultB != "2" {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with one int
	result, _ = runParserWithStr(`{ "a": 1 }`)
	if resultA2,_ := result.GetInt("a"); resultA2 != 1 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with multiple ints
	result, _ = runParserWithStr(`{ "a": 1, "b": 2 }`)
	resultAInt, _ := result.GetInt("a")
	resultBInt, _ := result.GetInt("b")
	if resultAInt != 1 || resultBInt != 2 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with one float
	result, _ = runParserWithStr(`{ "a": 1.11 }`)
	if resultFloat, _ := result.GetFloat("a"); resultFloat != 1.11 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with multiple floats
	result, _ = runParserWithStr(`{ "a": 1.11, "b": 2.22 }`)
	resultMultiFloatA, _ := result.GetFloat("a")
	resultMultiFloatB, _ := result.GetFloat("b")
	if resultMultiFloatA != 1.11 || resultMultiFloatB != 2.22 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}

	// Test valid JSON with nested object
	result, _ = runParserWithStr(`{ "a": { "b": { "c": 1 } } }`)
	objA, _ := result.GetObject("a")
	objB, _ := objA.GetObject("b")
	c, _ := objB.GetInt("c")
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
	arr, _ := result.GetArray("best_games")
	noita, _ := arr[0].GetInt("rank")
	smash, _ := arr[1].GetInt("rank")
	if noita != 1 || smash != 2 {
		t.Error(fmt.Sprintf("Unexpected parsed result: %s", result))
	}
}
