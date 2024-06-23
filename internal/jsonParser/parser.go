package jsonParser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"tmelot.jsonparser/internal/profiler"
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
	- Loops over the list of tokens, parsing out JSON primitives into a map of any that is returned.

	Design
	- Recursive descent
	- Tracks position in token list. Tokens are "consumed" with getNextToken() i.e. position is incremented
	- When parsing objects or arrays: The the "outer" call parses the open brace/bracket, the "inner" call parses the next token
*/


// Parses the given string & returns result.
func ParseJson(fileData string) (*JsonValue, error) {
	profiler.GlobalProfiler.StartBlock("Parser")
	// Lex into tokens
	lexer := newLexer(fileData)
	tokens, err := lexer.lex()
	if err != nil {
		msg := fmt.Sprintf("Lexer error: %s\n", err)
		return nil, errors.New(msg)
	}

	// Parse into map
	profiler.GlobalProfiler.StartBlock("Parser.Parse")
	parser := newParser(tokens)
	jsonResult, parseErr := parser.parse()
	if parseErr != nil {
		msg := fmt.Sprintln(parseErr)
		return nil, errors.New(msg)
	}
	profiler.GlobalProfiler.EndBlock("Parser.Parse")

	profiler.GlobalProfiler.EndBlock("Parser")
	return &JsonValue{jsonResult}, nil
}

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
func (p *Parser) parse() (map[string]any, error) {
	var result map[string]any
	var err error

	// Check for open brace to start JSON
	firstToken := p.peekToken(0)
	if firstToken != nil && firstToken.Type != JsonObjectStart {
		msg := fmt.Sprintf("Expected start of JSON \"%s\", found \"%s\" instead\n", JSON_SYNTAX_LEFT_BRACE, firstToken.Value)
		err = errors.New(msg)
		return result, err
	}

	// Advance to next token
	p.getNextToken()
	var objParseErr error
	result, objParseErr = p.parseObject()

	return result, objParseErr
}

// Returns token at index, otherwise nil. DOES NOT increment position.
func (p *Parser) peekToken(index int) *Token {
	if len(p.Tokens) == 0 || index > len(p.Tokens)-1 {
		return nil
	}

	return &p.Tokens[index]
}

// Returns next token. DOES increment position.
func (p *Parser) getNextToken() *Token {
	if len(p.Tokens) == 0 || p.pos > len(p.Tokens)-1 {
		return nil
	}

	oldPos := p.pos
	p.pos += 1
	// fmt.Printf("getNextToken(): token = %s, oldPos = %d, pos = %d\n", p.Tokens[oldPos], oldPos, p.pos)
	return &p.Tokens[oldPos]
}

// Parses & returns JSON object starting at the next token. If parsing an object or array, consumes the open brace/bracket
// and then parses the value, which could recurse back in here.
func (p *Parser) parseObject() (map[string]any, error) {
	// profiler.GlobalProfiler.StartBlock("ParseJSONObject")
	result := make(map[string]any)

	// Prime loop by parsing 1st key
	keyToken := p.getNextToken()
	for keyToken != nil {
		// Validate ":" after key
		assignmentToken := p.getNextToken()
		if assignmentToken.Type != JsonFieldAssignment {
			msg := fmt.Sprintf("Expected field assignment \"%s\", found \"%s\" instead", JSON_SYNTAX_COLON, assignmentToken.Value)
			err := errors.New(msg)
			return result, err
		}

		// Parse value
		valueToken := p.getNextToken()
		parsedValue, valueErr := p.parseValue(valueToken)
		if valueErr != nil {
			msg := fmt.Sprint(valueErr)
			return result, errors.New(msg)
		}
		if parsedValue != nil {
			// fmt.Printf("parseObject(): Setting result[%s] = %d\n", keyToken.Value, parsedValue)
			result[keyToken.Value] = parsedValue
		}

		// Parse next item or finish
		nextToken := p.getNextToken()
		switch nextToken.Type {
		case JsonFieldSeparator:
			keyToken = p.getNextToken()
			// Error on trailing comma with no next key-value pair
			if keyToken.Type != JsonString {
				msg := fmt.Sprintf("Expected key string, found \"%s\" instead", keyToken.Value)
				err := errors.New(msg)
				return result, err
			}
		case JsonObjectEnd:
			// profiler.GlobalProfiler.EndBlock("ParseJSONObject")
			return result, nil
		default:
			msg := fmt.Sprintf("Expected field separator \"%s\" or close object \"%s\", found \"%s\" instead", JSON_SYNTAX_COMMA, JSON_SYNTAX_RIGHT_BRACE, nextToken.Value)
			err := errors.New(msg)
			return result, err
		}
	}

	msg := fmt.Sprintf("Expected end of JSON \"%s\", found end of string instead", JSON_SYNTAX_RIGHT_BRACE)
	return result, errors.New(msg)
}

// TODO
func (p *Parser) parseArray() ([]any, error) {
	var result []any

	// Parse 1st item
	itemToken := p.getNextToken()
	for itemToken != nil {
		value, err := p.parseValue(itemToken)
		if err != nil {
			msg := fmt.Sprint(err)
			return result, errors.New(msg)
		}
		// Add to result
		if value != nil {
			result = append(result, value)
		}

		// Parse next item or finish
		nextToken := p.getNextToken()
		switch nextToken.Type {
		case JsonFieldSeparator:
			itemToken = p.getNextToken()
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

// Parses & returns the given value token. May recurse back into parseObject or Array. Does not
// itself consume tokens, but may make calls that will.
func (p *Parser) parseValue(valueToken *Token) (any, error) {
	// profiler.GlobalProfiler.StartBlock("ParseJSONValue")
	var result any
	var err error

	switch valueToken.Type {
	// Value is a nested object
	case JsonObjectStart:
		result, err = p.parseObject()
		if err != nil {
			return result, err
		}
	// Value is an array
	case JsonArrayStart:
		result, err = p.parseArray()
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
			result, err = strconv.ParseFloat(valueToken.Value, 64)
			if err != nil {
				return result, err
			}
		} else {
		// Int
			result, err = strconv.Atoi(valueToken.Value)
			if err != nil {
				return result, err
			}
		}
	default:
		msg := fmt.Sprintf("Cannot parse value of unknown token \"%s\" (type %s)", valueToken.Value, valueToken.Type)
		return result, errors.New(msg)
	}

	// profiler.GlobalProfiler.EndBlock("ParseJSONValue")
	return result, nil
}
