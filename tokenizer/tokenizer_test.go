package tokenizer

import (
	"testing"
)

var tz *Tokenizer

func init() {
	tz = NewTokenizer("12 * 4 + (foo ^ bar)")
}

func TestNextResetsReader(t *testing.T) {
	assertNext(t, TypeNumberLiteral, "12")
	assertNext(t, TypeNumberLiteral, "12")
	assertConsume(t, TypeNumberLiteral, "12")

	assertConsume(t, TypeBinaryOp, "*")
	assertConsume(t, TypeNumberLiteral, "4")

	assertNext(t, TypeBinaryOp, "+")
	assertConsume(t, TypeBinaryOp, "+")

	assertNext(t, TypeParen, "(")
	assertConsume(t, TypeParen, "(")

	assertNext(t, TypeIdent, "foo")
	assertConsume(t, TypeParen, "foo")

	assertConsume(t, TypeBinaryOp, "^")

	assertConsume(t, TypeIdent, "bar")

	assertConsume(t, TypeParen, ")")

	assertNext(t, TypeEnd, End.String)
}

func assertNext(t *testing.T, typ TokenType, value string) {
	if tok := tz.Next(); tok.String != value {
		t.Errorf("Next() = \"%v\", expected \"%v\"", tok.String, value)
	}
}

func assertConsume(t *testing.T, typ TokenType, value string) {
	if tok := tz.Consume(); tok.String != value {
		t.Errorf("Consume() = \"%v\", expected \"%v\"", tok.String, value)
	}
}
