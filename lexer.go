package main

import (
	"errors"
	"fmt"
)

/*
	Lexer works by scanning thru a JSON string & splitting it into tokens. It keeps track
	of the current scan position (see NOTE-1 below for more on that).
*/

// JSON syntax
const JSON_SYNTAX_WHITESPACE = " \n\t"
const JSON_SYNTAX_LEFT_BRACE = "{"
const JSON_SYNTAX_RIGHT_BRACE = "}"
const JSON_SYNTAX_LEFT_BRACKET = "["
const JSON_SYNTAX_RIGHT_BRACKET = "]"
const JSON_SYNTAX_COLON = ":"
const JSON_SYNTAX_COMMA = ","
const JSON_SYNTAX_QUOTE = "\""

// Identifies which type of JSON syntax the token represents
type TokenType string
const (
	JsonObjectStart TokenType = "ObjectStart"
	JsonObjectEnd TokenType = "ObjectEnd"
	JsonArrayStart TokenType = "ArrayStart"
	JsonArrayEnd TokenType = "ArrayEnd"
	JsonFieldAssignment TokenType = "FieldAssignment"
	JsonFieldSeparator TokenType = "FieldSeparator"
	JsonString TokenType = "String"
	JsonNumber TokenType = "Number"
	// TODO: True, False, Null
)

// Represents a lexed token
type Token struct {
	Type	TokenType
	Value 	string
}

type Lexer struct {
	Debug 	 bool
	data 	 string
	pos 	 int
}

// Create & return a new Lexer instance
func newLexer(data string) *Lexer {
	return &Lexer{
		data: data,
	}
}

// Returns string of data starting at pos, which represents the unlexed data
func (l *Lexer) getUnlexedData() string {
	return l.data[l.pos:]
}

// Lexes the internal string
// NOTE-1: This is the ONLY function that advances the lexer's position. The functions that do
// the actual lexing only return the number of characters consumed, which lex() uses to advance
// the position.
func (l *Lexer) lex() ([]Token, error) {
	var tokens []Token

	for l.pos < len(l.data) {
		l.DebugPrintf("pos = %d, len = %d\n", l.pos, len(l.data))

		// Lex JSON syntax
		syntaxToken, syntaxCharsRead := l.lexJsonSyntax()
		if syntaxCharsRead > 0 {
			tokens = append(tokens, *syntaxToken)
			l.pos += syntaxCharsRead
			continue
		}

		// Lex whitespace (which just ignores it)
		wsCharsRead := l.lexJsonWhitespace()
		if wsCharsRead > 0 {
			l.pos += wsCharsRead
			continue
		}

		// Lex strings
		stringToken, stringCharsRead, err := l.lexString()
		if err != nil {
			return tokens, err
		}
		if stringCharsRead > 0 {
			tokens = append(tokens, *stringToken)
			l.pos += stringCharsRead
			continue
		}

		// Lex numbers
		// NOTE: Numbers are read as strings. Later the parser will convert to correct data type.
		numberToken, numberCharsRead := l.lexNumber()
		if numberCharsRead > 0 {
			tokens = append(tokens, *numberToken)
			l.pos += numberCharsRead
			continue
		}

		// TODO: Lex bools
		// TODO: Lex null

		err = errors.New(fmt.Sprintf("Unexpected character \"%s\"", l.getUnlexedData()))
		return tokens, err
	}

	return tokens, nil
}

// Scans for consecutive whitespace & returns number of characters consumed.
// (Whitespace is thrown away)
func (l *Lexer) lexJsonWhitespace() int {
	numCharsRead := 0

	for _,s := range l.getUnlexedData() {
		foundWhitespace := false
		l.DebugPrintf("Scanning char %c for whitespace\n", s)
		for _,ws := range JSON_SYNTAX_WHITESPACE {
			if s == ws {
				l.DebugPrintf("	Match! %c is whitespace\n", s)
				numCharsRead += 1
				foundWhitespace = true
				break
			}
		}

		if foundWhitespace == false {
			break
		}
	}

	l.DebugPrintf("Consumed %d characters of whitespace\n", numCharsRead)
	return numCharsRead
}

// Scans for JSON syntax & returns it with number of characters consumed.
func (l *Lexer) lexJsonSyntax() (*Token, int) {
	s := string(l.getUnlexedData()[0])
	l.DebugPrintf("Checking %s for JSON syntax...", s)

	found := true // Weird default but it cuts verbosity in cases below
	var tokenType TokenType

	switch s {
	case JSON_SYNTAX_LEFT_BRACE:
		tokenType = JsonObjectStart
	case JSON_SYNTAX_RIGHT_BRACE:
		tokenType = JsonObjectEnd
	case JSON_SYNTAX_LEFT_BRACKET:
		tokenType = JsonArrayStart
	case JSON_SYNTAX_RIGHT_BRACKET:
		tokenType = JsonArrayEnd
	case JSON_SYNTAX_COLON:
		tokenType = JsonFieldAssignment
	case JSON_SYNTAX_COMMA:
		tokenType = JsonFieldSeparator
	default:
		found = false
	}

	if found {
		token := Token{Type: tokenType, Value: s}
		return &token, 1
	} else {
		return nil, 0
	}
}

// Scans for strings (like "a_string") & returns it along with number of characters consumed.
func (l *Lexer) lexString() (*Token, int, error) {
	s := l.getUnlexedData()
	lexedStr := ""
	numCharsRead := 0

	// Read past starting quote
	if string(s[0]) == JSON_SYNTAX_QUOTE {
		s = s[1:]
		numCharsRead += 1
	} else {
		l.DebugPrintf("%s is not a string\n", string(s[0]))
		return nil, numCharsRead, nil
	}

	// Scan string until we find closing quote
	for _,c := range s {
		l.DebugPrintf("Checking %c for string...\n", c)
		if string(c) == JSON_SYNTAX_QUOTE {
			numCharsRead += 1
			l.DebugPrintf("Returning lexed string %s\n", lexedStr)
			return &Token{Type: JsonString, Value: lexedStr}, numCharsRead, nil
		} else {
			numCharsRead += 1
			lexedStr += string(c)
			l.DebugPrintln("lexedStr =", lexedStr)
		}
	}

	// Error becasue we ran off edge of string without finding end quote
	err := errors.New(fmt.Sprint("End quote for string not found"))
	return nil, numCharsRead, err
}

// Scans for numbers (like "1" or "1.234") & returns it along with number of characters consumed.
func (l *Lexer) lexNumber() (*Token, int) {
	s := l.getUnlexedData()
	lexedStr := ""
	numCharsRead := 0

	for _,c := range s {
		l.DebugPrintf("Checking %c for number... ", c)
		isDigit := c >= '0' && c <= '9'
		isSymbol := c == '-' || c == '.'
		if isDigit || isSymbol {
			lexedStr += string(c)
			numCharsRead += 1
		} else {
			l.DebugPrintln("Found non number character, returning")
			break
		}
	}

	if numCharsRead > 0 {
		return &Token{Type: JsonNumber, Value: lexedStr}, numCharsRead
	} else {
		return nil, 0
	}
}

// TODO: Scans for bools (true or false) & returns it along with number of characters consumed.
// func (l *Lexer) lexBool() (*Token, int) {

func (l *Lexer) DebugPrintf(format string, a ...interface{}) {
	if l.Debug {
		fmt.Printf(format, a...)
	}
}

func (l *Lexer) DebugPrintln(a ...interface{}) {
	if l.Debug {
		fmt.Println(a...)
	}
}
