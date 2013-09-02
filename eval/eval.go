package eval

import (
	"fmt"
	"github.com/chlu/parser-experiments/ast"
	"math"
)

type Context struct {
	Variables map[string]interface{}
}

func NewContext() *Context {
	v := make(map[string]interface{})
	return &Context{v}
}

func Evaluate(n ast.Node, c *Context) (interface{}, error) {
	switch i := n.(type) {
	case *ast.NodeIdent:
		return c.Variables[string(*i)], nil
	case *ast.NodeNumber:
		return int(*i), nil
	case *ast.NodeBinaryOp:
		switch i.Op {
		case ast.OpBinPlus, ast.OpBinMinus, ast.OpBinMult, ast.OpBinDiv, ast.OpBinExp:
			v1, err := Evaluate(i.Lft, c)
			if err != nil {
				return nil, err
			}
			i1, isNumber := v1.(int)
			if !isNumber {
				return nil, fmt.Errorf("Multiplication is only supported on numbers, got %v", v1)
			}
			v2, err := Evaluate(i.Rgt, c)
			if err != nil {
				return nil, err
			}
			i2, isNumber := v2.(int)
			if !isNumber {
				return nil, fmt.Errorf("Multiplication is only supported on numbers, got %v", v2)
			}
			switch i.Op {
			case ast.OpBinPlus:
				return i1 + i2, nil
			case ast.OpBinMinus:
				return i1 - i2, nil
			case ast.OpBinMult:
				return i1 * i2, nil
			case ast.OpBinDiv:
				return i1 / i2, nil
			case ast.OpBinExp:
				return int(math.Pow(float64(i1), float64(i2))), nil
			default:
				panic("This should never happen")
			}
		case ast.OpAnd, ast.OpOr:
			v1, err := Evaluate(i.Lft, c)
			if err != nil {
				return nil, err
			}
			i1, isBool := v1.(bool)
			if !isBool {
				return nil, fmt.Errorf("Boolean operators are only supported on boolean expressions, got %v", v1)
			}
			v2, err := Evaluate(i.Rgt, c)
			if err != nil {
				return nil, err
			}
			i2, isBool := v2.(bool)
			if !isBool {
				return nil, fmt.Errorf("Boolean operators are only supported on boolean expressions, got %v", v2)
			}
			switch i.Op {
			case ast.OpAnd:
				return i1 && i2, nil
			case ast.OpOr:
				return i1 || i2, nil
			default:
				panic("This should never happen")
			}
		case ast.OpEquals:
			v1, err := Evaluate(i.Lft, c)
			if err != nil {
				return nil, err
			}
			v2, err := Evaluate(i.Rgt, c)
			if err != nil {
				return nil, err
			}
			return v1 == v2, nil
		case ast.OpUnaMinus:
			v1, err := Evaluate(i.Lft, c)
			if err != nil {
				return nil, err
			}
			i1, isNumber := v1.(int)
			if !isNumber {
				return nil, fmt.Errorf("Unary minus is only supported on numbers, got %v", v1)
			}
			return -i1, nil
		}
		return nil, fmt.Errorf("Cannot evaluate operation %v", i.Op)
	}
	return nil, fmt.Errorf("Cannot evaluate node %v", n)
}
