package main

import (
	"errors"
	"fmt"
)


type Parser struct {
	Tokens  []Token
	pos 	int
}

func newParser(tokens []Token) *Parser {
	return &Parser{
		Tokens: tokens,
	}
}

func (p *Parser) GetNextToken() *Token {
	if p.pos < len(p.Tokens) {
		token := p.Tokens[p.pos]
		p.pos += 1
		return &token
	} else {
		return nil
	}
}

func (p *Parser) IsFirstToken() bool {
	return p.pos == 0
}

func (p *Parser) Parse() (string, error) {
	result := ""

	// Prime & loop over tokens
	t := p.GetNextToken()
	for t != nil {
		fmt.Printf("Token %d: %s\n", p.pos, t)

		if p.IsFirstToken() && t.Type != JsonObjectStart {
			msg := fmt.Sprintf("Unexpected start of JSON, found %s, expected %s", t.Value, JSON_SYNTAX_LEFT_BRACE)
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
