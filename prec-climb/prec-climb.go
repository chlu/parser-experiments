package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chlu/parser-experiments/prec-climb/parser"
)

var debug = false

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "No expression given")
		os.Exit(1)
	}

	exp := flag.Arg(0)

	p := parser.NewParser(exp)

	fmt.Println(exp, ":", p.Parse())
}
