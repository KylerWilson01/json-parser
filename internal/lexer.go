package internal

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

const (
	ILLEGAL          TokenType = "Illegal"
	EOF              TokenType = "eof"
	OPENING_CURLY    TokenType = "{"
	CLOSING_CURLY    TokenType = "}"
	OPENING_BRACKET  TokenType = "["
	CLOSEING_BRACKET TokenType = "]"
	COLON            TokenType = ":"
	STRING           TokenType = "string"
	NUMBER           TokenType = "number"
)

func NewLexer(input string) *Lexer {
	l := Lexer{input: input}
	l.readChar()
	return &l
}

func (l *Lexer) NextToken() Token {
	var t Token

	l.skipWhiteSpace()

	switch l.ch {
	case '{':
		t = Token{Literal: string(l.ch), Type: OPENING_CURLY}
	case '}':
		t = Token{Literal: string(l.ch), Type: CLOSING_CURLY}
	case ':':
		t = Token{Literal: string(l.ch), Type: COLON}
	case '"':
		t = l.readString()
	case 0:
		t = Token{Literal: "", Type: EOF}
	default:
		if l.isNumber(l.ch) || l.ch == '-' {
			t = l.readNumber()
		} else {
			t = Token{Literal: "", Type: ILLEGAL}
		}
	}

	l.readChar()

	return t
}

func (l *Lexer) isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) readNumber() Token {
	position := l.position

	for l.isNumber(l.peek()) || l.peek() == '.' || l.peek() == '-' {
		l.readChar()
	}

	return Token{Type: NUMBER, Literal: l.input[position:l.readPosition]}
}

func (l *Lexer) readString() Token {
	var t Token
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' {
			t = Token{Type: STRING, Literal: l.input[position:l.position]}
			break
		}
		if l.ch == 0 {
			if l.input[l.position-1] != '"' {
				t = Token{Type: ILLEGAL, Literal: l.input[position:l.position]}
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

func (l *Lexer) prev() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition-1]
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
