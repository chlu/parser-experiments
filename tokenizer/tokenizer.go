package tokenizer

import (
	"fmt"
	"io"
	"strings"
)

type Token struct {
	String string
	TokenType
}

type TokenType uint8

var (
	End        = Token{"<<end>>", TypeEnd}
	ParenOpen  = Token{"(", TypeEnd}
	ParenClose = Token{")", TypeEnd}
)

const (
	TypeEnd TokenType = iota
	TypeIdent
	TypeNumberLiteral
	TypeBinaryOp
	TypeParen
)

type Tokenizer struct {
	reader io.ByteScanner
}

// Create a new Tokenizer for the given expression
func NewTokenizer(exp string) *Tokenizer {
	t := &Tokenizer{reader: strings.NewReader(exp)}
	return t
}

// Peek at the next token from expression without changing the position
func (t *Tokenizer) Next() Token {
	return t.scanToken(true)
}

// Consume the next token and change position to next token
func (t *Tokenizer) Consume() Token {
	return t.scanToken(false)
}

func (t *Tokenizer) scanToken(reset bool) Token {
	var (
		b   byte
		err error
	)
	for {
		b, err = t.reader.ReadByte()
		if err != nil {
			return End
		}
		// If we have whitespace, just skip over it
		if b != ' ' {
			break
		}
	}
	if reset {
		defer t.reader.UnreadByte()
	}

	if isDigit(b) {
		s := string(b) + readWhile(t.reader, isDigit, reset)
		return Token{s, TypeNumberLiteral}
	}

	if isIdent(b) {
		s := string(b) + readWhile(t.reader, isIdent, reset)
		return Token{s, TypeIdent}
	}

	if isBinaryOp(b) {
		return Token{string(b), TypeBinaryOp}
	}

	if b == '(' {
		return ParenOpen
	}
	if b == ')' {
		return ParenClose
	}

	panic(fmt.Sprintf("Unexpected token '%v'", string(b)))
}

// ---- Utility functions

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isIdent(b byte) bool {
	return 'a' <= b && b <= 'z'
}

func isBinaryOp(b byte) bool {
	return b == '+' || b == '-' || b == '*' || b == '/' || b == '^' || b == '='
}

// Read from reader while f holds true, unread bytes if reset == true
func readWhile(reader io.ByteScanner, f func(b byte) bool, reset bool) string {
	s := ""

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return s
		} else {
			if !f(b) {
				reader.UnreadByte()
				return s
			} else {
				if reset {
					defer reader.UnreadByte()
				}
				s += string(b)
			}
		}
	}
}
