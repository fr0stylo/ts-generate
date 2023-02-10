package lexer

import (
	"io"
	"unicode"

	"github.com/samber/lo"
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

var tokens = []byte{
	byte(Token_OpenSquareBracket),
	byte(Token_CloseSquareBracket),
	byte(Token_OpenCurlyBracket),
	byte(Token_CloseCurlyBracket),
	byte(Token_Quote),
	byte(Token_Comma),
	byte(Token_Colon),
	byte(Token_NewLine),
}

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
	Text string
	Type TokenType
}

func (t *Token) String() string {
	return t.Text
}

// func Tokenize(reader io.Reader) ([]*Token, error) {
// 	tokens := []*Token{}
// 	l := NewLexer(reader)

// 	for {
// 		t, err := l.Read()
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			return nil, err
// 		}
// 		tokens = append(tokens, t)
// 	}

// 	return tokens, nil
// }

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
		return Token{Type: EOF, Text: "EOF"}, io.EOF
	}
	switch b {
	case byte(Token_OpenSquareBracket):
		return Token{Text: string(b), Type: Array_OpenBracket}, nil
	case byte(Token_CloseSquareBracket):
		return Token{Text: string(b), Type: Array_CloseBracket}, nil
	case byte(Token_OpenCurlyBracket):
		return Token{Text: string(b), Type: Object_OpenBracket}, nil
	case byte(Token_CloseCurlyBracket):
		return Token{Text: string(b), Type: Object_CloseBracket}, nil
	case byte(Token_Colon):
		return Token{Text: string(b), Type: Colon}, nil
	case byte(Token_Comma):
		return Token{Text: string(b), Type: Comma}, nil
	case byte(Token_Quote):
		return Token{Text: string(b), Type: Quote}, nil
	default:
		sym := []byte{b}
		for nb, err := l.read(); err == nil && lo.IndexOf(tokens, nb) == -1; nb, err = l.read() {
			sym = append(sym, nb)
		}
		l.unread()

		return Token{Text: string(sym), Type: Symbol}, nil
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
