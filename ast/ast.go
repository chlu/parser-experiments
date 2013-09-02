package ast

import (
	"fmt"
	"github.com/chlu/parser-experiments/tokenizer"
	"strconv"
)

type Node interface {
	String() string
}

type NodeIdent string
type NodeNumber int
type NodeBinaryOp struct {
	Op  Operator
	Lft Node
	Rgt Node
}

func MakeIdent(tok tokenizer.Token) *NodeIdent {
	i := NodeIdent(tok.String)
	return &i
}

func MakeNumber(tok tokenizer.Token) *NodeNumber {
	i, err := strconv.ParseInt(tok.String, 0, 0)
	if err != nil {
		panic("Could not convert string to number")
	}
	n := NodeNumber(i)
	return &n
}

func MakeNode(op Operator, lft, rgt Node) *NodeBinaryOp {
	return &NodeBinaryOp{op, lft, rgt}
}

func (n *NodeIdent) String() string {
	return string(*n)
}

func (n *NodeNumber) String() string {
	return strconv.Itoa(int(*n))
}

func (n *NodeBinaryOp) String() string {
	if n.Rgt != nil {
		return fmt.Sprintf("%v(%v,%v)", n.Op.String(), n.Lft.String(), n.Rgt.String())
	}
	return fmt.Sprintf("%v(%v)", n.Op.String(), n.Lft.String())
}

type Operator uint8

const (
	OpOr Operator = iota
	OpAnd
	OpEquals
	OpBinPlus
	OpBinMinus
	OpUnaMinus
	OpBinMult
	OpBinDiv
	OpBinExp
)

func ToBinaryOp(tok tokenizer.Token) Operator {
	switch tok.String {
	case "+":
		return OpBinPlus
	case "-":
		return OpBinMinus
	case "*":
		return OpBinMult
	case "/":
		return OpBinDiv
	case "^":
		return OpBinExp
	case "=":
		return OpEquals
	}
	panic(fmt.Sprintf("Unknown binary op token %s", tok))
}

func IsUnary(tok tokenizer.Token) bool {
	return tok.String == "-"
}

func ToUnaryOp(tok tokenizer.Token) Operator {
	switch tok.String {
	case "-":
		return OpUnaMinus
	}
	panic(fmt.Sprintf("Unknown binary op token %s", tok))
}

func (op Operator) String() string {
	switch op {
	case OpBinPlus:
		return "+"
	case OpBinMinus:
		return "-"
	case OpBinMult:
		return "*"
	case OpBinDiv:
		return "/"
	case OpBinExp:
		return "^"
	case OpUnaMinus:
		return "-"
	case OpEquals:
		return "="
	}
	panic(fmt.Sprintf("Unknown operator %v", op))
}

type Associativity uint8

const (
	AssocLeft Associativity = iota
	AssocRight
)

func (op Operator) Prec() int {
	switch op {
	case OpOr:
		return 0
	case OpAnd:
		return 1
	case OpEquals:
		return 2
	case OpBinPlus, OpBinMinus:
		return 3
	case OpUnaMinus:
		return 4
	case OpBinMult, OpBinDiv:
		return 5
	case OpBinExp:
		return 6
	}
	panic(fmt.Sprintf("Undefined prec for operator %s", op.String()))
}

func (op Operator) Associativity() Associativity {
	switch op {
	case OpBinExp:
		return AssocRight
	}
	return AssocLeft
}
