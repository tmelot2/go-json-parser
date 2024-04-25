package main

import (
	"errors"
	"fmt"
)

/*
	TODO: Parser docs
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

// Retuns token at given position, or nil if does not exist. Does NOT advance position.
func (p *Parser) peekToken(pos int) *Token {
	if pos < len(p.Tokens) - 1 {
		token := p.Tokens[pos]
		return &token
	} else {
		return nil
	}
}

// Gets token at next position & returns it. DOES advance position.
// TODO: Need this?
func (p *Parser) getNextToken() *Token {
	if p.pos < len(p.Tokens) - 1 {
		p.pos += 1
		token := p.Tokens[p.pos]
		fmt.Printf("Returning next token \"%s\" (type %s)\n", token.Value, token.Type)
		return &token
	} else {
		return nil
	}
}

// // Advances position by the given int, caps at token len - 1.
// func (p *Parser) advancePosition(num int) {
// 	p.pos += num
// 	if p.pos > len(p.Tokens) - 1 {
// 		p.pos = len(p.Tokens) - 1
// 	}
// }

func (p *Parser) isFirstToken() bool {
	return p.pos == 0
}

func (p *Parser) isLastToken() bool {
	return p.pos == len(p.Tokens)-1
}

func (p *Parser) Parse(isRoot bool) (map[string]string, error) {
	// result := ""

	// // Prime & loop over tokens
	// t := p.peekToken(0)
	// for t != nil {
	// 	// Error on missing open bracket to start JSON
	// 	if p.isFirstToken() && t.Type != JsonObjectStart {
	// 		msg := fmt.Sprintf(`Unexpected start of JSON, found "%s", expected "%s"`, t.Value, JSON_SYNTAX_LEFT_BRACE)
	// 	} else if p.isLastToken() && t.Type != JsonObjectEnd {
	// 	// Error on missing end bracket to end JSON
	// 		msg := fmt.Sprintf(`Unexpected end of JSON, found "%s", expected "%s"`, t.Value, JSON_SYNTAX_RIGHT_BRACE)
	// 		return "", errors.New(msg)
	// 	} else {
	// 		result += fmt.Sprintf("%s ", t.Value)
	// 	}

	// 	t = p.getNextToken()
	// }

	// return result, nil

	var t *Token

	if isRoot {
		t = p.peekToken(0)
	} else {
		t = p.getNextToken()
	}

	if t != nil {
		fmt.Printf("Parse(): Checking token %s, pos = %d\n", t, p.pos)

		if isRoot && t.Type != JsonObjectStart {
			msg := fmt.Sprintf(`Unexpected start of JSON, found "%s", expected "%s"`, t.Value, JSON_SYNTAX_LEFT_BRACE)
			nilResult := make(map[string]string)
			return nilResult, errors.New(msg)
		}
	}

	// Consume the token & parse
	result, err := p.parseObject()

	return result, err
}

func (p *Parser) parseObject() (map[string]string, error) {
	result := make(map[string]string)

	t := p.getNextToken()
	for t != nil {
		fmt.Printf("parseObject(): Checking token %s, pos = %d\n", t, p.pos)
		if t.Type == JsonObjectEnd {
			return result, nil
		}

		// Validate & parse key
		if t.Type != JsonString {
			msg := fmt.Sprintf("Expected %s type for object key, found %s instead", JsonString, t.Type)
			return result, errors.New(msg)
		}
		key := t.Value

		// Validate field separator
		if p.getNextToken().Type != JsonFieldAssignment {
			msg := fmt.Sprintf(`Expected field assignment "%s" after key, found "%s" (type %s) instead`, JSON_SYNTAX_COLON, t.Value, t.Type)
			return result, errors.New(msg)
		}

		// Parse value
		valueToken := p.getNextToken()
		value := valueToken.Value

		// Set result
		result[key] = value

		// Parse end or next token
		nextToken := p.getNextToken()
		if nextToken.Type == JsonObjectEnd {
			return result, nil
		} else if nextToken.Type != JsonFieldSeparator {
			msg := fmt.Sprintf(`Expected field separator "%s" after key, found "%s" instead`, JSON_SYNTAX_COMMA, nextToken.Value)
			return result, errors.New(msg)
		}

		t = p.getNextToken()
	}

	// Error on missing close bracket
	if t == nil {
		msg := fmt.Sprintf(`Expected end of object "%s", found end of file instead`, JSON_SYNTAX_RIGHT_BRACE)
		return result, errors.New(msg)
	} else {
		msg := fmt.Sprintf(`Expected field assignment "%s" after key, found "%s" (type %s) instead`, JSON_SYNTAX_COLON, t.Value, t.Type)
		return result, errors.New(msg)
	}
}

func (p *Parser) parseArray() string {
	return ""
}
