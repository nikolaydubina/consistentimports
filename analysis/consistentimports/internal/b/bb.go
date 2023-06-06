package b

import (
	"fmt"
	Str2 "strings"

	aN2 "golang.org/x/tools/go/analysis"
)

func AA() {
	fmt.Println(Str2.Join([]string{"a", "b", "c"}, ","))

	var x *aN2.Analyzer
	if x == nil {
		fmt.Println("x is nil")
	}
}
