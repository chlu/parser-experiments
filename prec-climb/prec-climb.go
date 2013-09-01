package main

import (
	"fmt"
	"os"

	"github.com/chlu/parser-experiments/prec-climb/parser"
)

var debug = false

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "No expression given")
		os.Exit(1)
	}

	exp := args[1]

	p := parser.NewParser()
	p.Debug = true

	n, err := p.Parse(exp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing expression: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(exp, ":", n)
}
