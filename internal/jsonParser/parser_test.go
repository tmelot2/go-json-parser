package jsonParser

/*
	Tests the parser & JsonValue classes.
*/

import (
	// "fmt"
	"testing"

	"tmelot.jsonparser/internal/assert"
)

func runParserWithStr(s string) (*JsonValue, error) {
	result, err := ParseJson(s)
	return result, err
}

func TestParserInvalidJson(t *testing.T) {
	// Test empty string
	_, err := runParserWithStr("")
	assert.NotNil(t, err, "Expected error for blank JSON, did not error")

	// Test missing start JSON
	_, err = runParserWithStr("}")
	assert.NotNil(t, err, "Expected error on missing start JSON open brace, did not error")

	// Test missing start JSON
	_, err = runParserWithStr(`"hello": "world"}`)
	assert.NotNil(t, err, "Expected error on missing start JSON open brace, did not error")

	// Test missing end JSON
	_, err = runParserWithStr("{")
	assert.NotNil(t, err, "Expected error on missing end JSON close brace, did not error")

	// Test missing field assignment ":"
	_, err = runParserWithStr(`{"hello" "world"}`)
	assert.NotNil(t, err, "Expected error on missing field assignment \":\", did not error")

	// Test missing field separator ","
	_, err = runParserWithStr(`{ "a": 1 "b": 2 }`)
	assert.NotNil(t, err, "Expected error on missing field separator \",\", did not error")

	// Test missing quote on key
	_, err = runParserWithStr(`{ a: 1 }`)
	assert.NotNil(t, err, "Expected error on missing quotes around key, did not error")

	// Test invalid object trailing comma
	_, err = runParserWithStr(`{ "a": 1, "b": 2, }`)
	assert.NotNil(t, err, "Expected error on trailing comma, did not error")

	// Test invalid array trailing comma
	_, err = runParserWithStr(`{ "a": [1,2,] }`)
	assert.NotNil(t, err, "Expected error on trailing array comma, did not error")

	// Test invalid array trailing commaa
	_, err = runParserWithStr(`{ "a": [1,2,,] }`)
	assert.NotNil(t, err, "Expected error on multiple trailing array commas, did not error")
}

func TestParserValidJson(t *testing.T) {
	// Test valid JSON with 1 string
	result, _ := runParserWithStr(`{ "a": "1" }`)
	val, _ := result.GetString("a")
	assert.Equal(t, val, "1")

	// Test valid JSON with multiple strings
	result, _ = runParserWithStr(`{ "a": "1", "b": "2" }`)
	resultA, _ := result.GetString("a")
	resultB, _ := result.GetString("b")
	assert.Equal(t, resultA, "1")
	assert.Equal(t, resultB, "2")

	// Test valid JSON with one int
	result, _ = runParserWithStr(`{ "a": 1 }`)
	resultA2, _ := result.GetInt("a")
	assert.Equal(t, resultA2, 1)

	// Test valid JSON with multiple ints
	result, _ = runParserWithStr(`{ "a": 1, "b": 2 }`)
	resultAInt, _ := result.GetInt("a")
	resultBInt, _ := result.GetInt("b")
	assert.Equal(t, resultAInt, 1)
	assert.Equal(t, resultBInt, 2)

	// Test valid JSON with one float
	result, _ = runParserWithStr(`{ "a": 1.11 }`)
	resultFloat, _ := result.GetFloat("a")
	assert.Equal(t, resultFloat, 1.11)

	// Test valid JSON with multiple floats
	result, _ = runParserWithStr(`{ "a": 1.11, "b": 2.22 }`)
	resultMultiFloatA, _ := result.GetFloat("a")
	resultMultiFloatB, _ := result.GetFloat("b")
	assert.Equal(t, resultMultiFloatA, 1.11)
	assert.Equal(t, resultMultiFloatB, 2.22)

	// Test valid JSON with nested object
	result, _ = runParserWithStr(`{ "a": { "b": { "c": 1 } } }`)
	objA, _ := result.GetObject("a")
	objB, _ := objA.GetObject("b")
	c, _ := objB.GetInt("c")
	assert.Equal(t, c, 1)

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
	assert.Equal(t, noita, 1)
	assert.Equal(t, smash, 2)

	// Test JSON with all field types
	result, _ = runParserWithStr(`{
		"str": "string",
		"int": 1,
		"obj": {
			"str": "obj.string",
			"int": 12
		},
		"arr": [11,22,33],
		"true": true,
		"false": false
	}`)
	strVal, _ := result.GetString("str")
	intVal, _ := result.GetInt("int")
	assert.Equal(t, strVal, "string")
	assert.Equal(t, intVal, 1)

	objVal, _ := result.GetObject("obj")
	strVal, _ = objVal.GetString("str")
	intVal, _ = objVal.GetInt("int")
	assert.Equal(t, strVal, "obj.string")
	assert.Equal(t, intVal, 12)

	arrVal, _ := result.GetArray("arr")
	arr0, _ := arrVal[0].GetInt("")
	arr1, _ := arrVal[1].GetInt("")
	arr2, _ := arrVal[2].GetInt("")
	assert.Equal(t, arr0, 11)
	assert.Equal(t, arr1, 22)
	assert.Equal(t, arr2, 33)

	trueVal, _ := result.GetBool("true")
	falseVal, _ := result.GetBool("false")
	assert.Equal(t, trueVal, true)
	assert.Equal(t, falseVal, false)

	assert.Finished()
}
