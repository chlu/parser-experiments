package parser

import (
	"testing"
)

var p *Parser

func init() {
	p = NewParser()
}

func TestParseExp(t *testing.T) {
	exp := "a * b - c * d - e * f = g * h - i * j - k * l"
	out := "=(-(-(*(a,b),*(c,d)),*(e,f)),-(-(*(g,h),*(i,j)),*(k,l)))"

	n, err := p.Parse(exp)
	if err != nil {
		t.Fatalf("Error parsing expression \"%s\": %v", exp, err)
	}
	if s := n.String(); s != out {
		t.Fatalf("Parsed expression \"%s\" to %s but expected %s", exp, s, out)
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := p.Parse("")
	if err == nil {
		t.Fatalf("Expected error for empty expression")
	}
}
