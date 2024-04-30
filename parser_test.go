package main

import (
	// "testing"
)


func runParserWithStr(s string) (map[string]interface{}, error) {
	// lexer := newLexer(s)
	// // lexer.Debug = true
	// tokens, _ := lexer.lex()
	// parser := newParser(tokens)
	// // parser.Debug = true
	// result, err := parser.Parse()
	// return result, err

	result, err := ParseJson(s)
	return result, err
}
