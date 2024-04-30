package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
	TODO: Unsafe dynamic typing docs
	- Trading off type safety to gain simplicity (but is it, really?)
	    - Parsed JSON is a nested map of varying types FOR EVERY FIELD
	    - (I think) you MUST cast each map traversal to the correct type in calling code, which will be very messy
	    	- Maybe there's a way to metaprogram/reflect out of this problem with a type definition
	    - BUT we get a parser that is extremely simple & has zero reflection
	    - Alternate idea: What about using reflection to map using builtin struct JSON markup?
	    	- seems like a ton of work to do this, def not simple

	    "my current design is trading off type safety to gain simplicity. the way it's implemented, my parsing functions will return a string or a number or a map. the calling client code will be responsible for safely traversing the map. since i know the exact use case, and it's only 1 JSON format, that seems like a fair trade off."

	Parser
	- Loops over the list of tokens, parsing out JSON primitives into a map that is returned.

	Design
	- Recursive descent
	- Tracks position in token list. Tokens are "consumed" with GetNextToken() i.e. position is incremented
	- When parsing objects or arrays: The the "outer" call parses the open brace/bracket, the "inner" call parses the next token
*/


type Parser struct {
	Debug 	 bool
	Tokens  []Token
	pos 	int
}

func newParser(tokens []Token) *Parser {
	return &Parser{
		Tokens: tokens,
	}
}

// Parses tokens & returns map of result, or a partial result with an error. It tries to
// return as much as it's parsed so far.
func (p *Parser) Parse() (map[string]any, error) {
	var result map[string]any
	var err error

	// Check for open brace to start JSON
	firstToken := p.PeekToken(0)
	if firstToken != nil && firstToken.Type != JsonObjectStart {
		msg := fmt.Sprintf("Expected start of JSON \"%s\", found \"%s\" instead\n", JSON_SYNTAX_LEFT_BRACE, firstToken.Value)
		err = errors.New(msg)
		return result, err
	}

	// Advance to next token
	p.GetNextToken()
	var objParseErr error
	result, objParseErr = p.ParseObject()

	return result, objParseErr
}

// Returns token at index, otherwise nil. DOES NOT increment position.
func (p *Parser) PeekToken(index int) *Token {
	if len(p.Tokens) == 0 || index > len(p.Tokens)-1 {
		return nil
	}

	return &p.Tokens[index]
}

// Returns next token. DOES increment position.
func (p *Parser) GetNextToken() *Token {
	if len(p.Tokens) == 0 || p.pos > len(p.Tokens)-1 {
		return nil
	}

	oldPos := p.pos
	p.pos += 1
	// fmt.Printf("GetNextToken(): token = %s, oldPos = %d, pos = %d\n", p.Tokens[oldPos], oldPos, p.pos)
	return &p.Tokens[oldPos]
}

// Parses & returns JSON object starting at the next token. If parsing an object or array, consumes the open brace/bracket
// and then parses the value, which could recurse back in here.
func (p *Parser) ParseObject() (map[string]any, error) {
	result := make(map[string]any)

	// Prime loop by parsing 1st key
	keyToken := p.GetNextToken()
	for keyToken != nil {
		// Validate ":" after key
		assignmentToken := p.GetNextToken()
		if assignmentToken.Type != JsonFieldAssignment {
			msg := fmt.Sprintf("Expected field assignment \"%s\", found \"%s\" instead", JSON_SYNTAX_COLON, assignmentToken.Value)
			err := errors.New(msg)
			return result, err
		}

		// Parse value
		valueToken := p.GetNextToken()
		parsedValue, valueErr := p.ParseValue(valueToken)
		if valueErr != nil {
			msg := fmt.Sprintf("Error: %s", valueErr)
			return result, errors.New(msg)
		}
		if parsedValue != nil {
			// fmt.Printf("ParseObject(): Setting result[%s] = %d\n", keyToken.Value, parsedValue)
			result[keyToken.Value] = parsedValue
		}

		// Parse next item or finish
		nextToken := p.GetNextToken()
		switch nextToken.Type {
		case JsonFieldSeparator:
			keyToken = p.GetNextToken()
		case JsonObjectEnd:
			return result, nil
		default:
			msg := fmt.Sprintf("Expected field separator \"%s\" or close object \"%s\", found \"%s\" instead", JSON_SYNTAX_COMMA, JSON_SYNTAX_RIGHT_BRACE, nextToken.Value)
			err := errors.New(msg)
			return result, err
		}
	}

	return result, nil
}

// TODO
func (p *Parser) ParseArray() ([]any, error) {
	var result []any

	// Parse 1st item
	itemToken := p.GetNextToken()
	for itemToken != nil {
		value, err := p.ParseValue(itemToken)
		if err != nil {
			msg := fmt.Sprintf("Error: %s", err)
			return result, errors.New(msg)
		}
		// Add to result
		if value != nil {
			result = append(result, value)
		}

		// Parse next item or finish
		nextToken := p.GetNextToken()
		switch nextToken.Type {
		case JsonFieldSeparator:
			itemToken = p.GetNextToken()
		case JsonArrayEnd:
			return result, nil
		// TODO: Objects or other arrays
		default:
			msg := fmt.Sprintf("Expected field separator \"%s\" or close object \"%s\", found \"%s\" instead", JSON_SYNTAX_COMMA, JSON_SYNTAX_RIGHT_BRACKET, itemToken.Value)
			return result, errors.New(msg)
		}
	}

	return result, nil
}

// Parses & returns the given value token. May recurse back into ParseObject or Array. Does not
// itself consume tokens, but may make calls that will.
func (p *Parser) ParseValue(valueToken *Token) (any, error) {
	var result any

	switch valueToken.Type {
	// Value is a nested object
	case JsonObjectStart:
		var err error
		result, err = p.ParseObject()
		if err != nil {
			return result, err
		}
	// Value is an array
	case JsonArrayStart:
		var err error
		result, err = p.ParseArray()
		if err != nil {
			return result, err
		}
	// Value is a string
	case JsonString:
		result = valueToken.Value
	// Value is a number
	case JsonNumber:
		// TODO: How to handle strconv errors?
		// Float
		if strings.Contains(valueToken.Value, ".") {
			result, _ = strconv.ParseFloat(valueToken.Value, 64)
		} else {
		// Int
			result, _ = strconv.Atoi(valueToken.Value)
		}
	}

	return result, nil
}
