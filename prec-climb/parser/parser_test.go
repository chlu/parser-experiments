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
		t.Errorf("Error parsing expression \"%s\": %v", err)
	}
	if s := n.String(); s != out {
		t.Errorf("Parsed expression \"%s\" to %s but expected %s", exp, s, out)
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := p.Parse("")
	if err == nil {
		t.Errorf("Expected error for empty expression")
	}
}
