package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
	TODO: Parser docs
	TODO: Unsafe dynamic typing docs
	- Trading off type safety to gain simplicity (but is it, really?)
	    - Parsed JSON is a nested map of varying types FOR EVERY FIELD
	    - (I think) you MUST cast each map traversal to the correct type in calling code, which will be very messy
	    	- Maybe there's a way to metaprogram/reflect out of this problem with a type definition
	    - BUT we get a parser that is extremely simple & has zero reflection
	    - Alternate idea: What about using reflection to map using builtin struct JSON markup?
	    	- seems like a ton of work to do this, def not simple

	    "my current design is trading off type safety to gain simplicity. the way it's implemented, my parsing functions will return a string or a number or a map. the calling client code will be responsible for safely traversing the map. since i know the exact use case, and it's only 1 JSON format, that seems like a fair trade off."
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

func (p *Parser) PeekToken(index int) *Token {
	if len(p.Tokens) == 0 || index > len(p.Tokens)-1 {
		return nil
	}

	return &p.Tokens[index]
}

func (p *Parser) GetNextToken() *Token {
	if len(p.Tokens) == 0 || p.pos > len(p.Tokens)-1 {
		return nil
	}

	oldPos := p.pos
	p.pos += 1
	fmt.Printf("GetNextToken(): token = %s, oldPos = %d, pos = %d\n", p.Tokens[oldPos], oldPos, p.pos)
	return &p.Tokens[oldPos]
}

func (p *Parser) ParseObject() (map[string]any, error) {
	fmt.Println("	parsing new object")
	result := make(map[string]any)

	keyToken := p.GetNextToken()
	for keyToken != nil {
		// Validate : after key
		fmt.Println("	validating :")
		assignmentToken := p.GetNextToken()
		if assignmentToken.Type != JsonFieldAssignment {
			msg := fmt.Sprintf("Expected field assignment \"%s\", found \"%s\" instead", JSON_SYNTAX_COLON, assignmentToken.Value)
			err := errors.New(msg)
			return result, err
		}

		// Parse value
		fmt.Println("	parsing value")
		valueToken := p.GetNextToken()
		// Value is a nested object
		if valueToken.Type == JsonObjectStart {
			// p.GetNextToken() // Throw away open bracket
			var valueErr error
			result[keyToken.Value], valueErr = p.ParseObject()
			if valueErr != nil {
				return result, valueErr
			}
		} else if valueToken.Type == JsonString {
		// Value is a string
			result[keyToken.Value] = valueToken.Value
		} else if valueToken.Type == JsonNumber {
		// Value is a number
			// Float
			if strings.Contains(valueToken.Value, ".") {
				result[keyToken.Value], _ = strconv.ParseFloat(valueToken.Value, 64)
			} else {
			// Int
				result[keyToken.Value], _ = strconv.Atoi(valueToken.Value)
			}
		}

		// Parse next or finish
		fmt.Println("	parsing next or finishing")
		nextToken := p.GetNextToken()
		if nextToken.Type == JsonFieldSeparator {
			keyToken = p.GetNextToken()
		} else if nextToken.Type == JsonObjectEnd {
			return result, nil
		} else {
			msg := fmt.Sprintf("Expected field separator \"%s\" or close object \"%s\", found \"%s\" instead", JSON_SYNTAX_COMMA, JSON_SYNTAX_RIGHT_BRACE, nextToken.Value)
			err := errors.New(msg)
			return result, err
		}
	}

	return result, nil
}