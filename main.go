package main

import (
	"github.com/nikolaydubina/consistentimports/analysis/consistentimports"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(consistentimports.Analyzer) }
