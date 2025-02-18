package internal

import "fmt"

// TokenError holds the error for when a token is illegal
type TokenError struct {
	msg, arg string
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("%s %s", t.msg, t.arg)
}

// TokenType is a string.
type TokenType string

// Token holds what a token should represent.
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer is what we use to make sure that all Tokens are valid.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	Tokens       []Token
}

const (
	// Illegal shows that the token is not valid
	Illegal TokenType = "Illegal"

	// OpeningCurly is what shows the start of an object
	OpeningCurly TokenType = "{"
	// ClosingCurly is what shows the end of an object
	ClosingCurly TokenType = "}"
	// OpeningBracket marks the begging of an array
	OpeningBracket TokenType = "["
	// CloseingBracket marks the end of an array
	CloseingBracket TokenType = "]"

	// Colon seperates the key and value
	Colon TokenType = ":"
	// Comma seperates the values
	Comma TokenType = ","

	// Null marks a primitive null
	Null TokenType = "null"
	// String marks a primitive string
	String TokenType = "string"
	// Number marks a primitive number
	Number TokenType = "number"
	// True marks a primitive true boolean
	True TokenType = "true"
	// False marks a primitive false boolean
	False TokenType = "false"
)

// NewLexer creates a pointer to a Lexer.
func NewLexer(input string) *Lexer {
	l := Lexer{input: input}
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
			l.Tokens = append(l.Tokens, Token{Literal: string(l.ch), Type: OpeningCurly})
		case '}':
			if l.Tokens[idx-1].Type == Comma {
				return fmt.Errorf(
					"Character \"%s\" is not valid in the given position",
					l.Tokens[idx-1].Literal,
				)
			}
			l.Tokens = append(l.Tokens, Token{Literal: string(l.ch), Type: ClosingCurly})
		case ':':
			l.Tokens = append(l.Tokens, Token{Literal: string(l.ch), Type: Colon})
		case ',':
			l.Tokens = append(l.Tokens, Token{Literal: string(l.ch), Type: Comma})
		case '"':
			l.Tokens = append(l.Tokens, l.readString())
		case 0:
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
				return fmt.Errorf("Character \"%s\" is not valid", string(l.ch))
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
		t = Token{Type: True, Literal: string(True)}
	case 'f':
		for _, c := range False[1:] {
			if c != rune(l.peek()) {
				return nil, &TokenError{arg: string(l.peek()), msg: "Character was not false"}
			}
			l.readChar()
		}
		t = Token{Type: False, Literal: string(False)}
	case 'n':
		for _, c := range Null[1:] {
			if c != rune(l.peek()) {
				return nil, &TokenError{arg: string(l.peek()), msg: "Character was not null"}
			}
			l.readChar()
		}
		t = Token{Type: Null, Literal: string(Null)}
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

	for l.isNumber(l.peek()) || l.peek() == '.' || l.peek() == '-' {
		l.readChar()
	}

	return Token{Type: Number, Literal: l.input[position:l.readPosition]}
}

func (l *Lexer) readString() Token {
	var t Token
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' {
			t = Token{Type: String, Literal: l.input[position:l.position]}
			break
		}
		if l.ch == 0 {
			if l.input[l.position-1] != '"' {
				t = Token{Type: Illegal, Literal: l.input[position:l.position]}
			}
			break
		}
	}

	return t
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

func (l *Lexer) isWhiteSpace(ch1, ch2 byte) bool {
	return ch1 == ' ' || (ch1 == '\\' && (ch2 == 't' || ch2 == 'n' || ch2 == 'r'))
}

func (l *Lexer) skipWhiteSpace() {
	for {
		if l.ch == ' ' {
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
