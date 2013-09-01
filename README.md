Parser experiments
==================

Some parser experiments in Go (golang) based on the precedence climbing paper from Theodore Norvell (http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm).

Installation
------------

    go install github.com/chlu/parser-experiments/prec-climb

You'll need a working Go environment for this.

Usage
-----

Parsing expressions:

    $GOPATH/bin/prec-climb "x + y * (a / 4)"
