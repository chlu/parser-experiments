package eval

import (
	"github.com/chlu/parser-experiments/ast"
	"testing"
)

func TestEvaluateIdent(t *testing.T) {
	n := ast.NodeIdent("foo")

	c := NewContext()
	c.Variables["foo"] = 42

	out, err := Evaluate(&n, c)
	if err != nil {
		t.Fatalf("Error evaluating node: %v", err)
	}

	if out != 42 {
		t.Fatalf("Expected \"foo\" == 42, got %v", out)
	}
}

func TestEvaluateNumber(t *testing.T) {
	n := ast.NodeNumber(42)

	c := NewContext()

	out, err := Evaluate(&n, c)
	if err != nil {
		t.Fatalf("Error evaluating node: %v", err)
	}

	if out != 42 {
		t.Fatalf("Expected \"42\" == 42, got %v", out)
	}
}

func TestEvaluateCalculation(t *testing.T) {
	n0 := ast.NodeNumber(2)
	n1 := ast.NodeNumber(21)
	n := ast.NodeBinaryOp{ast.OpBinMult, &n0, &n1}

	c := NewContext()

	out, err := Evaluate(&n, c)
	if err != nil {
		t.Fatalf("Error evaluating node: %v", err)
	}

	if out != 42 {
		t.Fatalf("Expected \"2*21\" == 42, got %v", out)
	}

	n = ast.NodeBinaryOp{ast.OpBinExp, &n0, &n1}

	out, err = Evaluate(&n, c)
	if err != nil {
		t.Fatalf("Error evaluating node: %v", err)
	}

	if out != 2097152 {
		t.Fatalf("Expected \"2^21\" == 2097152, got %v", out)
	}
}

func TestEvaluateBoolean(t *testing.T) {
	n0 := ast.NodeNumber(4)
	n1 := ast.NodeNumber(16)
	n2 := ast.NodeBinaryOp{ast.OpBinMult, &n0, &n0}
	nt := ast.NodeBinaryOp{ast.OpEquals, &n2, &n1}
	nf := ast.NodeBinaryOp{ast.OpEquals, &n0, &n1}

	c := NewContext()

	out, err := Evaluate(&nt, c)
	if err != nil {
		t.Fatalf("Error evaluating node: %v", err)
	}

	if out != true {
		t.Fatalf("Expected \"4*4=16\" == true, got %v", out)
	}

	out, err = Evaluate(&nf, c)
	if err != nil {
		t.Fatalf("Error evaluating node: %v", err)
	}

	if out != false {
		t.Fatalf("Expected \"4=16\" == false, got %v", out)
	}

}

func TestEvaluateUnaryOp(t *testing.T) {
	n0 := ast.NodeNumber(42)
	n := ast.NodeBinaryOp{ast.OpUnaMinus, &n0, nil}

	c := NewContext()

	out, err := Evaluate(&n, c)
	if err != nil {
		t.Fatalf("Error evaluating node: %v", err)
	}

	if out != -42 {
		t.Fatalf("Expected \"-42\" == 42, got %v", out)
	}
}
