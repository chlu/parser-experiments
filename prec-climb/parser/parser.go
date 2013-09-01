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
	"github.com/chlu/parser-experiments/tokenizer"
	"strconv"
)

// ----

type node interface {
	String() string
}

type nodeIdent string
type nodeNumber int
type nodeBinaryOp struct {
	op  operator
	lft node
	rgt node
}

func mkIdent(tok tokenizer.Token) *nodeIdent {
	i := nodeIdent(tok.String)
	return &i
}

func mkNumber(tok tokenizer.Token) *nodeNumber {
	i, err := strconv.ParseInt(tok.String, 0, 0)
	if err != nil {
		panic("Could not convert string to number")
	}
	n := nodeNumber(i)
	return &n
}

func mkNode(op operator, lft, rgt node) *nodeBinaryOp {
	return &nodeBinaryOp{op, lft, rgt}
}

func (n *nodeIdent) String() string {
	return string(*n)
}

func (n *nodeNumber) String() string {
	return strconv.Itoa(int(*n))
}

func (n *nodeBinaryOp) String() string {
	return fmt.Sprintf("%v(%v,%v)", n.op.String(), n.lft.String(), n.rgt.String())
}

type operator uint8

const (
	opOr operator = iota
	opAnd
	opEquals
	opBinPlus
	opBinMinus
	opUnaMinus
	opBinMult
	opBinDiv
	opBinExp
)

type associativity uint8

const (
	assocLeft associativity = iota
	assocRight
)

func toBinaryOp(tok tokenizer.Token) operator {
	switch tok.String {
	case "+":
		return opBinPlus
	case "-":
		return opBinMinus
	case "*":
		return opBinMult
	case "/":
		return opBinDiv
	case "^":
		return opBinExp
	default:
		panic(fmt.Sprintf("Unknown binary op token %s", tok))
	}
}

func isUnary(tok tokenizer.Token) bool {
	return tok.String == "-"
}

func toUnaryOp(tok tokenizer.Token) operator {
	switch tok.String {
	case "-":
		return opUnaMinus
	default:
		panic(fmt.Sprintf("Unknown binary op token %s", tok))
	}
}

func (op operator) String() string {
	switch op {
	case opBinPlus:
		return "+"
	case opBinMinus:
		return "-"
	case opBinMult:
		return "*"
	case opBinDiv:
		return "/"
	case opBinExp:
		return "^"
	case opUnaMinus:
		return "-"
	}
	panic(fmt.Sprintf("Unknown operator %v", op))
}

func (op operator) prec() int {
	switch op {
	case opOr:
		return 0
	case opAnd:
		return 1
	case opEquals:
		return 2
	case opBinPlus, opBinMinus:
		return 3
	case opUnaMinus:
		return 4
	case opBinMult, opBinDiv:
		return 5
	case opBinExp:
		return 6
	default:
		panic("Undefined prec")
	}
}

func (op operator) associativity() associativity {
	switch op {
	case opBinExp:
		return assocRight
	default:
		return assocLeft
	}
}

// ----

type Parser struct {
	tz *tokenizer.Tokenizer
}

func NewParser(exp string) *Parser {
	p := &Parser{tz: tokenizer.NewTokenizer(exp)}
	return p
}

func (p *Parser) Parse() node {
	t := p.parse_Exp(0)
	p.expect(tokenizer.End)
	return t
}

func (p *Parser) parse_Exp(pp int) node {
	t := p.parse_P()
	for {
		if n := p.tz.Next(); n.TokenType == tokenizer.TypeBinaryOp && toBinaryOp(n).prec() >= pp {
			var q int
			op := toBinaryOp(n)
			p.tz.Consume()

			switch op.associativity() {
			case assocRight:
				q = op.prec()
			case assocLeft:
				q = 1 + op.prec()
			}
			t1 := p.parse_Exp(q)
			t = mkNode(op, t, t1)
		} else {
			break
		}
	}
	return t
}

func (p *Parser) parse_P() node {
	n := p.tz.Next()
	if isUnary(n) {
		op := toUnaryOp(n)
		p.tz.Consume()
		q := op.prec()
		t := p.parse_Exp(q)
		return mkNode(op, t, nil)
	} else if n.String == "(" {
		p.tz.Consume()
		t := p.parse_Exp(0)
		p.expect(tokenizer.ParenClose)
		return t
	} else if n.TokenType == tokenizer.TypeIdent {
		t := mkIdent(n)
		p.tz.Consume()
		return t
	} else if n.TokenType == tokenizer.TypeNumberLiteral {
		t := mkNumber(n)
		p.tz.Consume()
		return t
	} else {
		panic(fmt.Sprintf("Could not parse token %s", n.String))
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
