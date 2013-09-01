/*

	A parser for the following grammar G:

		E --> Exp(0)
		Exp(p) --> P {B Exp(q)}
		P --> U Exp(q) | "(" E ")" | v
		B --> "+" | "-"  | "*" |"/" | "^" | "||" | "&&" | "="
		U --> "-"

*/
package parser

import (
	"fmt"
	"github.com/chlu/parser-experiments/ast"
	"github.com/chlu/parser-experiments/tokenizer"
	"strings"
)

type Parser struct {
	tz    *tokenizer.Tokenizer
	Debug bool
}

func NewParser() *Parser {
	p := &Parser{}
	return p
}

func (p *Parser) Parse(exp string) (ast.Node, error) {
	p.tz = tokenizer.NewTokenizer(exp)

	t, err := p.parse_Exp(0, 0)
	if err != nil {
		return nil, err
	}
	err = p.expect(tokenizer.End)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (p *Parser) parse_Exp(pp, lvl int) (ast.Node, error) {
	if p.Debug {
		fmt.Printf("%sparse_Exp(%d)\n", strings.Repeat("  ", lvl), pp)
	}
	t, err := p.parse_P(lvl)
	if err != nil {
		return nil, err
	}
	for {
		if n := p.tz.Next(); n.TokenType == tokenizer.TypeBinaryOp && ast.ToBinaryOp(n).Prec() >= pp {
			var q int
			op := ast.ToBinaryOp(n)
			p.tz.Consume()

			switch op.Associativity() {
			case ast.AssocRight:
				q = op.Prec()
			case ast.AssocLeft:
				q = 1 + op.Prec()
			}
			if p.Debug {
				fmt.Printf("%s- op %s\n", strings.Repeat("  ", lvl), op.String())
			}
			t1, err := p.parse_Exp(q, lvl+1)
			if err != nil {
				return nil, err
			}
			t = ast.MakeNode(op, t, t1)
			if p.Debug {
				fmt.Printf("%s-> %s\n", strings.Repeat("  ", lvl), t.String())
			}
		} else {
			break
		}
	}
	return t, nil
}

func (p *Parser) parse_P(lvl int) (ast.Node, error) {
	n := p.tz.Next()
	if ast.IsUnary(n) {
		op := ast.ToUnaryOp(n)
		p.tz.Consume()
		q := op.Prec()
		t, err := p.parse_Exp(q, lvl+1)
		if err != nil {
			return nil, err
		}
		t1 := ast.MakeNode(op, t, nil)
		if p.Debug {
			fmt.Printf("%s-> %s\n", strings.Repeat("  ", lvl), t1.String())
		}
		return t1, nil
	} else if n.String == "(" {
		p.tz.Consume()
		t, err := p.parse_Exp(0, lvl+1)
		if err != nil {
			return nil, err
		}
		err = p.expect(tokenizer.ParenClose)
		if err != nil {
			return nil, err
		}
		return t, nil
	} else if n.TokenType == tokenizer.TypeIdent {
		t := ast.MakeIdent(n)
		if p.Debug {
			fmt.Printf("%s-> %s\n", strings.Repeat("  ", lvl), t.String())
		}
		p.tz.Consume()
		return t, nil
	} else if n.TokenType == tokenizer.TypeNumberLiteral {
		t := ast.MakeNumber(n)
		if p.Debug {
			fmt.Printf("%s-> %s\n", strings.Repeat("  ", lvl), t.String())
		}
		p.tz.Consume()
		return t, nil
	} else {
		return nil, fmt.Errorf("Expected one of -, (, <ident> or <number> got %s", n.String)
	}
}

// ----

func (p *Parser) expect(tok tokenizer.Token) error {
	if ntok := p.tz.Next(); ntok == tok {
		p.tz.Consume()
		return nil
	} else {
		return fmt.Errorf("Expected token %s, but got %s", tok.String, ntok.String)
	}
}
