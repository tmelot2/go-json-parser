package main

import (
	"errors"
	"fmt"
)


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
func (p *Parser) PeekToken(pos int) *Token {
	if pos < len(p.Tokens) - 1 {
		token := p.Tokens[pos]
		return &token
	} else {
		return nil
	}
}

// Gets token at next position & returns it. DOES advance position.
func (p *Parser) GetNextToken() *Token {
	if p.pos < len(p.Tokens) - 1 {
		p.pos += 1
		token := p.Tokens[p.pos]
		return &token
	} else {
		return nil
	}
}

func (p *Parser) IsFirstToken() bool {
	return p.pos == 0
}

func (p *Parser) IsLastToken() bool {
	return p.pos == len(p.Tokens)-1
}

func (p *Parser) Parse() (string, error) {
	result := ""

	// Prime & loop over tokens
	t := p.PeekToken(0)
	for t != nil {
		fmt.Printf("wwwwwwwwww pos = %d, t = %s, manual t = %s\n", p.pos, t, p.Tokens[p.pos])
		// Error on missing open bracket to start JSON
		if p.IsFirstToken() && t.Type != JsonObjectStart {
			msg := fmt.Sprintf(`Unexpected start of JSON, found "%s", expected "%s"`, t.Value, JSON_SYNTAX_LEFT_BRACE)
			return "", errors.New(msg)
		} else if p.IsLastToken() && t.Type != JsonObjectEnd {
		// Error on missing end bracket to end JSON
			msg := fmt.Sprintf(`Unexpected end of JSON, found "%s", expected "%s"`, t.Value, JSON_SYNTAX_RIGHT_BRACE)
			return "", errors.New(msg)
		} else {
			result += fmt.Sprintf("%s ", t.Value)
		}

		t = p.GetNextToken()
	}

	return result, nil
}

func (p *Parser) ParseObject() string {
	return ""
}

func (p *Parser) ParseArray() string {
	return ""
}

func (p *Parser) ParseString() string {
	return ""
}

func (p *Parser) ParseNumber() string {
	return ""
}
