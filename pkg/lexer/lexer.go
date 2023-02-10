package lexer

import (
	"io"
	"unicode"
)

const (
	Token_OpenSquareBracket  JsonToken = '['
	Token_CloseSquareBracket JsonToken = ']'
	Token_OpenCurlyBracket   JsonToken = '{'
	Token_CloseCurlyBracket  JsonToken = '}'
	Token_Quote              JsonToken = '"'
	Token_Comma              JsonToken = ','
	Token_Colon              JsonToken = ':'
	Token_NewLine            JsonToken = '\n'
)

const (
	Array_OpenBracket   = iota
	Array_CloseBracket  = iota
	Object_OpenBracket  = iota
	Object_CloseBracket = iota
	Comma               = iota
	Colon               = iota
	Symbol              = iota
	Quote               = iota
	EOF                 = iota
)

type JsonToken rune
type TokenType int

type Token struct {
	Text []byte
	Type TokenType
}

func (t *Token) String() string {
	return string(t.Text)
}

func Tokenize(reader io.Reader) ([]Token, error) {
	tokens := []Token{}
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	l := NewLexer(buf)

	for token, err := l.Read(); err == nil; token, err = l.Read() {
		tokens = append(tokens, token)
	}

	return tokens, nil
}

type Lexer struct {
	r   []byte
	pos int
}

func NewLexer(r []byte) *Lexer {
	return &Lexer{
		r:   r,
		pos: 0,
	}
}

func (l *Lexer) Read() (Token, error) {
	if l.pos >= len(l.r) {
		return Token{Type: EOF}, io.EOF
	}

	for unicode.IsSpace(rune(l.r[l.pos])) {
		l.pos++
	}

	b, err := l.read()
	if err != nil {
		return Token{Type: EOF, Text: nil}, io.EOF
	}
	switch b {
	case byte(Token_OpenSquareBracket):
		return Token{Text: []byte{b}, Type: Array_OpenBracket}, nil
	case byte(Token_CloseSquareBracket):
		return Token{Text: []byte{b}, Type: Array_CloseBracket}, nil
	case byte(Token_OpenCurlyBracket):
		return Token{Text: []byte{b}, Type: Object_OpenBracket}, nil
	case byte(Token_CloseCurlyBracket):
		return Token{Text: []byte{b}, Type: Object_CloseBracket}, nil
	case byte(Token_Colon):
		return Token{Text: []byte{b}, Type: Colon}, nil
	case byte(Token_Comma):
		return Token{Text: []byte{b}, Type: Comma}, nil
	case byte(Token_Quote):
		return Token{Text: []byte{b}, Type: Quote}, nil
	default:
		return l.parseSymbolic()
	}
}

func (l *Lexer) parseSymbolic() (Token, error) {
	sym := []byte{}
	l.unread()

	for nb, err := l.read(); err == nil && !isTerminationToken(nb); nb, err = l.read() {
		sym = append(sym, nb)
	}
	l.unread()

	return Token{Text: sym, Type: Symbol}, nil
}

func isTerminationToken(b byte) bool {
	switch b {
	case byte(Token_OpenSquareBracket),
		byte(Token_CloseSquareBracket),
		byte(Token_OpenCurlyBracket),
		byte(Token_CloseCurlyBracket),
		byte(Token_Quote),
		byte(Token_Comma),
		byte(Token_Colon),
		byte(Token_NewLine):
		return true
	default:
		return false
	}
}

func (l *Lexer) read() (byte, error) {
	if l.pos >= len(l.r) {
		return 0, io.EOF
	}
	b := l.r[l.pos]
	l.pos++
	return b, nil
}

func (l *Lexer) unread() {
	l.pos--
}

func (l *Lexer) Reset() {
	l.pos = 0
}
