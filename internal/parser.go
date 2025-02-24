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
	s := NewStack[TokenType]()
	for i, t := range p.tokens {
		switch t.Type {
		case OpeningCurly:
			if i == 0 {
				s.Push(OpeningCurly)
				continue
			}
			pt := p.tokens[i-1]
			if pt.Type != Colon &&
				(t.State == InsideArray && pt.Type != Comma && pt.Type != OpeningBracket) {
				return false, fmt.Errorf(
					"NAME_SEPARATOR should preceed LEFT_CURLY_BRACKET if this is a nested object or else nothing should preceed it, instead got: %v",
					pt.Literal,
				)
			}
			s.Push(OpeningCurly)
		case ClosingCurly:
			prevTkn := p.tokens[i-1].Type
			if prevTkn != ValueString && prevTkn != Number && prevTkn != True && prevTkn != False &&
				prevTkn != Null &&
				prevTkn != ClosingCurly &&
				prevTkn != CloseingBracket &&
				prevTkn != OpeningCurly {
				return false, fmt.Errorf(
					"VALUE_STRING or RIGHT_CURLY_BRACKET or LITERAL or RIGHT_SQUARE_BRACKET should preceed RIGHT_CURLY_BRACKET, instead got: %v",
					prevTkn,
				)
			}
			s.Pop()
		case OpeningBracket:
			if i == 0 {
				s.Push(OpeningBracket)
				continue
			}
			pt := p.tokens[i-1]
			if pt.Type != Colon && pt.Type != OpeningBracket && pt.Type != Comma {
				return false, fmt.Errorf(
					"NAME_SEPARATOR or LEFT_SQUARE_BRACKET should preceed LEFT_SQUARE_BRACKET, instead got: %v",
					pt.Literal,
				)
			}
			s.Push(OpeningBracket)
		case CloseingBracket:
			prevTkn := p.tokens[i-1].Type
			if prevTkn != ValueString && prevTkn != Number && prevTkn != True && prevTkn != False &&
				prevTkn != Null &&
				prevTkn != ClosingCurly &&
				prevTkn != CloseingBracket &&
				prevTkn != OpeningBracket {
				return false, fmt.Errorf(
					"VALUE_STRING or RIGHT_CURLY_BRACKET or LITERAL or RIGHT_SQUARE_BRACKETshould preceed RIGHT_SQUARE_BRACKET, instead got: %v",
					prevTkn,
				)
			}
			s.Pop()
		case NameString:
			prevTkn := p.tokens[i-1]
			if prevTkn.Type != OpeningCurly && prevTkn.Type != Comma {
				return false, fmt.Errorf(
					"OpeningCurly or Comma should preceed String, instead got: %v",
					prevTkn,
				)
			}
		case ValueString:
			prevTkn := p.tokens[i-1]
			if prevTkn.Type != Colon &&
				(t.State == InsideArray && prevTkn.Type != OpeningBracket && prevTkn.Type != Comma) {
				return false, fmt.Errorf(
					"NAME_SEPARATOR or LEFT_SQUARE_BRACKET (when within an array) or VALUE_SEPARATOR (when within an array) should preceed VALUE_STRING, instead got: %v",
					prevTkn,
				)
			}
		case Colon:
			prevTkn := p.tokens[i-1].Type
			if prevTkn != NameString {
				return false, fmt.Errorf(
					"String should preceed NAME_SEPARATOR, instead got: %v",
					prevTkn,
				)
			}
		case Comma:
			prevTkn := p.tokens[i-1].Type
			prevTknState := p.tokens[i-1].State
			if prevTknState == Invalid {
				return false, fmt.Errorf(
					"VALUE_SEPARATOR should not come after OpeningBrace or  RIGHT_CURLY_BRACKET (when the object isn't nested, got: %v",
					prevTkn,
				)
			}

			if prevTkn != OpeningCurly && (prevTkn != ValueString && prevTknState == InsideArray) &&
				prevTkn != Number &&
				prevTkn != Null &&
				prevTkn != True &&
				prevTkn != False {
				return false, fmt.Errorf(
					"RIGHT_CURLY_BRACKET or VALUE_STRING or NUMBER or LITERAL must precede VALUE_SEPARATOR, got: %v",
					prevTkn,
				)
			}
		case Number:
			prevTkn := p.tokens[i-1].Type
			if prevTkn != Colon &&
				(t.State == InsideArray && prevTkn != Comma && prevTkn != OpeningBracket) {
				return false, fmt.Errorf(
					"NAME_SEPARATOR or VALUE_SEPARATOR (within an array) or LEFT_SQUARE_BRACKET (within an array) should preceed NUMBER, instead got: %v",
					prevTkn,
				)
			}
		case True, False, Null:
			prevTkn := p.tokens[i-1]
			if prevTkn.Type != Colon &&
				(t.State == InsideArray && prevTkn.Type != Comma && prevTkn.Type != OpeningBracket) {
				return false, fmt.Errorf(
					"NAME_SEPARATOR should preceed LITERAL, instead got: %v",
					prevTkn,
				)
			}
		default:
			return false, fmt.Errorf("illegal token")
		}
	}

	if !s.IsEmpty() {
		return false, fmt.Errorf("Stack is not empty")
	}

	return true, nil
}
