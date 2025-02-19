package internal

import "fmt"

type Parser struct {
	tokens []Token
}

func NewParser(t []Token) *Parser {
	p := &Parser{}
	p.tokens = t
	return p
}

func (p *Parser) ParseTokens() (bool, error) {
	var e error
	result := false
	s := NewStack[TokenState]()
	for i, t := range p.tokens {
		if t.State == &StartArray || t.State == &StartObject {
			result = false
			e = fmt.Errorf("Start without an end")
			s.Push(*t.State)
			continue
		}

		if s.Peek() == nil {
			return false, fmt.Errorf("Empty stack")
		}

		if *t.State == EndObject && *s.Pop() == StartObject {
			if p.tokens[i-1].Type == Comma {
				result = false
				return result, fmt.Errorf("Comma is in an invalid place")
			}
			e = nil
			result = true
			continue
		}

		if *t.State == EndArray && *s.Pop() == StartArray {
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
