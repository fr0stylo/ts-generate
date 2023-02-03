package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

const (
	Token_OpenSquareBracket  JsonToken = '['
	Token_CloseSquareBracket           = ']'
	Token_OpenCurlyBracket             = '{'
	Token_CloseCurlyBracket            = '}'
	Token_Quote                        = '"'
	Token_Comma                        = ','
	Token_Colon                        = ':'
	Token_NewLine                      = '\n'
)

const (
	Array_OpenBracket   = iota
	Array_CloseBracket  = iota
	Object_OpenBracket  = iota
	Object_CloseBracket = iota
	Comma               = iota
	Colon               = iota
	Text                = iota
	Numeric             = iota
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

func Tokenize(reader io.Reader) ([]*Token, error) {
	tokens := []*Token{}
	l := NewLexer(reader)

	for {
		t, err := l.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}

	// r := bufio.NewReader(reader)

	// for x, _, err := r.ReadRune(); err == nil; x, _, err = r.ReadRune() {
	// 	strBuf := bytes.NewBufferString("")
	// 	strBuf.WriteRune(x)

	// 	switch x {
	// 	case rune(Token_OpenSquareBracket):
	// 		tokens = append(tokens, Token{Type: Array_OpenBracket, Text: strBuf.String()})
	// 	case rune(Token_CloseSquareBracket):
	// 		tokens = append(tokens, Token{Type: Array_CloseBracket, Text: strBuf.String()})
	// 	case rune(Token_OpenCurlyBracket):
	// 		tokens = append(tokens, Token{Type: Object_OpenBracket, Text: strBuf.String()})
	// 	case rune(Token_CloseCurlyBracket):
	// 		tokens = append(tokens, Token{Type: Object_CloseBracket, Text: strBuf.String()})
	// 	case rune(Token_Comma):
	// 		tokens = append(tokens, Token{Type: Comma, Text: strBuf.String()})
	// 	case rune(Token_Colon):
	// 		tokens = append(tokens, Token{Type: Colon, Text: strBuf.String()})
	// 	case rune(Token_Quote):
	// 		for peeked, err := r.Peek(1); err == nil && peeked[0] != byte(Token_Quote); peeked, err = r.Peek(1) {
	// 			sr, _, _ := r.ReadRune()
	// 			strBuf.WriteRune(sr)
	// 		}

	// 		sr, _, _ := r.ReadRune()
	// 		strBuf.WriteRune(sr)

	// 		tokens = append(tokens, Token{Type: Text, Text: strBuf.String()})
	// 	default:
	// 		if unicode.IsDigit(x) {
	// 			for peeked, err := r.Peek(1); err == nil && peeked[0] != byte(Token_Comma) && peeked[0] != byte(Token_NewLine); peeked, err = r.Peek(1) {
	// 				sr, _, _ := r.ReadRune()
	// 				strBuf.WriteRune(sr)
	// 			}

	// 			tokens = append(tokens, Token{Type: Numeric, Text: strBuf.String()})
	// 		}
	// 	}
	// }

	return tokens, nil
}

type Lexer struct {
	r   *bufio.Reader
	buf *bytes.Buffer
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r:   bufio.NewReader(r),
		buf: bytes.NewBuffer([]byte{}),
	}
}

func (l *Lexer) Read() (*Token, error) {
	l.buf.Reset()
	x, _, err := l.r.ReadRune()
	if err != nil {
		return nil, err
	}
	l.buf.WriteRune(x)
	switch x {
	case rune(Token_OpenSquareBracket):
		return &Token{Type: Array_OpenBracket, Text: l.buf.String()}, nil
	case rune(Token_CloseSquareBracket):
		return &Token{Type: Array_CloseBracket, Text: l.buf.String()}, nil
	case rune(Token_OpenCurlyBracket):
		return &Token{Type: Object_OpenBracket, Text: l.buf.String()}, nil
	case rune(Token_CloseCurlyBracket):
		return &Token{Type: Object_CloseBracket, Text: l.buf.String()}, nil
	case rune(Token_Comma):
		return &Token{Type: Comma, Text: l.buf.String()}, nil
	case rune(Token_Colon):
		return &Token{Type: Colon, Text: l.buf.String()}, nil
	case rune(Token_Quote):
		a, _ := l.r.ReadSlice(byte(Token_Quote))
		// for peeked, err := l.r.Peek(1); err == nil && peeked[0] != byte(Token_Quote); peeked, err = l.r.Peek(1) {
		// 	sr, _, _ := l.r.ReadRune()
		// 	l.buf.WriteRune(sr)
		// }

		// sr, _, _ := l.r.ReadRune()
		l.buf.Write(a)

		return &Token{Type: Text, Text: l.buf.String()}, nil
	default:
		if unicode.IsDigit(x) {
			for peeked, err := l.r.Peek(1); err == nil && peeked[0] != byte(Token_Comma) && peeked[0] != byte(Token_NewLine); peeked, err = l.r.Peek(1) {
				sr, _, _ := l.r.ReadRune()
				l.buf.WriteRune(sr)
			}

			return &Token{Type: Numeric, Text: l.buf.String()}, nil
		}

		return l.Read()
	}
}
