package internal

import (
	"fmt"
)

// TokenError holds the error for when a token is illegal
type TokenError struct {
	msg, arg string
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("%s %s", t.msg, t.arg)
}

// TokenType is a string.
type TokenType string

// TokenState is a string.
type TokenState string

// Token holds what a token should represent.
type Token struct {
	Type    TokenType
	Literal string
	State   TokenState
}

// Lexer is what we use to make sure that all Tokens are valid.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	state        *Stack[TokenState]
	Tokens       []Token
}

const (
	// Invalid state
	Invalid TokenState = "Invalid"
	// StartObject state
	StartObject TokenState = "StartObject"
	// EndObject state
	EndObject TokenState = "EndObject"
	// StartArray state
	StartArray TokenState = "StartArray"
	// EndArray state
	EndArray TokenState = "EndArray"
	// InsideObject state
	InsideObject TokenState = "InsideObject"
	// InsideArray state
	InsideArray TokenState = "InsideArray"
)

const (
	// Illegal shows that the token is not valid
	Illegal TokenType = "Illegal"

	// OpeningCurly is what shows the start of an object
	OpeningCurly TokenType = "{"
	// ClosingCurly is what shows the end of an object
	ClosingCurly TokenType = "}"
	// OpeningBracket marks the begging of an array
	OpeningBracket TokenType = "["
	// ClosingBracket marks the end of an array
	ClosingBracket TokenType = "]"

	// Colon seperates the key and value
	Colon TokenType = ":"
	// Comma seperates the values
	Comma TokenType = ","

	// Null marks a primitive null
	Null TokenType = "null"
	// ValueString marks a primitive string
	ValueString TokenType = "value string"
	// NameString marks a primitive string
	NameString TokenType = "name string"
	// Number marks a primitive number
	Number TokenType = "number"
	// True marks a primitive true boolean
	True TokenType = "true"
	// False marks a primitive false boolean
	False TokenType = "false"
)

// NewLexer creates a pointer to a Lexer.
func NewLexer(input string) *Lexer {
	l := Lexer{input: input, state: NewStack[TokenState]()}
	l.readChar()
	return &l
}

// ValidateTokens returns the next token.
func (l *Lexer) ValidateTokens() error {
	idx := 0
	for ; idx < len(l.input); idx++ {
		l.skipWhiteSpace()

		switch l.ch {
		case '{':
			l.Tokens = append(
				l.Tokens,
				Token{
					Literal: string(l.ch),
					Type:    OpeningCurly,
					State:   StartObject,
				},
			)
			l.state.Push(InsideObject)
		case '}':
			if s := l.state.Pop(); s != InsideObject {
				return &TokenError{"Should be inside an object. Instead got", string(s)}
			}
			l.Tokens = append(
				l.Tokens,
				Token{Literal: string(l.ch), Type: ClosingCurly, State: EndObject},
			)
		case '[':
			l.Tokens = append(
				l.Tokens,
				Token{Literal: string(l.ch), Type: OpeningBracket, State: StartArray},
			)
			l.state.Push(InsideArray)
		case ']':
			if s := l.state.Pop(); s != InsideArray {
				return &TokenError{"Should be inside an object. Instead got", string(s)}
			}
			l.Tokens = append(
				l.Tokens,
				Token{Literal: string(l.ch), Type: ClosingBracket, State: EndArray},
			)
		case ':':
			l.Tokens = append(
				l.Tokens,
				Token{Literal: string(l.ch), Type: Colon, State: l.findState()},
			)
		case ',':
			l.Tokens = append(
				l.Tokens,
				Token{Literal: string(l.ch), Type: Comma, State: l.findState()},
			)
		case '"':
			l.Tokens = append(l.Tokens, l.readString())
		case 0:
			if len(l.state.state) != 0 {
				return fmt.Errorf(
					"Length of the state should be 0. Instead got %d",
					len(l.state.state),
				)
			}
			return nil
		default:
			if l.isNumber(l.ch) || l.ch == '-' {
				l.Tokens = append(l.Tokens, l.readNumber())
			} else if l.isLiteral(l.ch) {
				literal, err := l.readLiteral()
				if err != nil {
					return err
				}
				l.Tokens = append(l.Tokens, *literal)
			} else {
				return &TokenError{"Not a legal token", string(l.ch)}
			}
		}

		l.readChar()
	}

	return nil
}

func (l *Lexer) readLiteral() (*Token, error) {
	var t Token
	switch l.ch {
	case 't':
		for _, c := range True[1:] {
			if c != rune(l.peek()) {
				return nil, &TokenError{arg: string(l.peek()), msg: "Character was not true"}
			}
			l.readChar()
		}
		t = Token{Type: True, Literal: string(True), State: l.findState()}
	case 'f':
		for _, c := range False[1:] {
			if c != rune(l.peek()) {
				return nil, &TokenError{arg: string(l.peek()), msg: "Character was not false"}
			}
			l.readChar()
		}
		t = Token{Type: False, Literal: string(False), State: l.findState()}
	case 'n':
		for _, c := range Null[1:] {
			if c != rune(l.peek()) {
				return nil, &TokenError{arg: string(l.peek()), msg: "Character was not null"}
			}
			l.readChar()
		}
		t = Token{Type: Null, Literal: string(Null), State: l.findState()}
	}
	return &t, nil
}

func (l *Lexer) isLiteral(ch byte) bool {
	return ch == 't' || ch == 'f' || ch == 'n'
}

func (l *Lexer) isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) readNumber() Token {
	position := l.position

	for {
		if l.isNumber(l.peek()) || l.peek() == '.' || l.peek() == '-' || l.peek() == '+' {
			l.readChar()
		} else if l.peek() == 'e' || l.peek() == 'E' {
			l.readChar()
		} else {
			break
		}
	}

	return Token{Type: Number, Literal: l.input[position:l.readPosition], State: l.findState()}
}

func (l *Lexer) readString() Token {
	position := l.position + 1

	for {
		l.readChar()
		if l.prev() == '\\' {
			if l.ch == '"' || l.ch == '\\' || l.ch == '/' || l.ch == 'b' ||
				l.peek() == 'f' ||
				l.peek() == 'n' ||
				l.peek() == 'r' ||
				l.peek() == 't' {
				l.readChar()
				continue
			}

			return Token{Type: Illegal, Literal: l.input[position:l.position], State: Invalid}
		} else if l.ch == '"' {
			t := Token{Literal: l.input[position:l.position], State: l.findState()}
			if t.Type == Illegal {
				return t
			}

			if len(l.Tokens) == 0 {
				t.Type = ValueString
				return t
			}

			pt := l.Tokens[len(l.Tokens)-1]

			if pt.Type == OpeningCurly || (pt.Type == Comma && l.state.Peek() == InsideObject) {
				t.Type = NameString
			} else if l.state.Peek() == InsideArray && (pt.Type == Comma || pt.Type == OpeningBracket) {
				t.Type = ValueString
			} else if pt.Type == Colon {
				t.Type = ValueString
			}
			return t
		} else if l.ch == 0 {
			if l.input[l.position-1] != '"' {
				return Token{Type: Illegal, Literal: l.input[position:l.position], State: Invalid}
			}
		}
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) prev() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition-1]
}

func (l *Lexer) isWhiteSpace(ch1, ch2 byte) bool {
	return ch1 == ' ' || (ch1 == '\\' && (ch2 == 't' || ch2 == 'n' || ch2 == 'r'))
}

func (l *Lexer) skipWhiteSpace() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
			continue
		}

		if l.ch == '\\' && (l.peek() == 't' || l.peek() == 'n' || l.peek() == 'r') {
			l.readChar()
			l.readChar()
			continue
		}
		break
	}
}

func (l *Lexer) findState() TokenState {
	if l.state.IsEmpty() {
		return Invalid
	}

	return l.state.Peek()
}
