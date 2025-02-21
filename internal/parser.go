package internal

import "fmt"

// Parser is used to parse the given tokens
type Parser struct {
	tokens []Token
}

// NewParser creates a new parser
func NewParser(t []Token) *Parser {
	p := &Parser{}
	p.tokens = t
	return p
}

// ParseTokens loops through all the tokens to make sure it's valid
func (p *Parser) ParseTokens() (bool, error) {
	var e error
	result := false
	s := NewStack[TokenState]()
	for i, t := range p.tokens {
		if t.State == StartArray || t.State == StartObject {
			result = false
			e = fmt.Errorf("Start without an end")
			s.Push(t.State)
			continue
		}

		if t.State == EndObject && s.Pop() == StartObject {
			if p.tokens[i-1].Type == Comma {
				result = false
				return result, fmt.Errorf("Comma is in an invalid place")
			}
			e = nil
			result = true
			continue
		}

		if t.State == EndArray && s.Pop() == StartArray {
			if p.tokens[i-1].Type == Comma {
				result = false
				return result, fmt.Errorf("Comma is in an invalid place")
			}
			e = nil
			result = true
			continue
		}

		if t.Type == Illegal {
			result = false
			return result, fmt.Errorf("Illegal character")
		}
	}

	return result, e
}
