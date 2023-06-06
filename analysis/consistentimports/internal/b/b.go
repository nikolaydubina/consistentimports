package b

import (
	"fmt"
	Str "strings"

	aN "golang.org/x/tools/go/analysis"
)

func A() {
	fmt.Println(Str.Join([]string{"a", "b", "c"}, ","))

	var x *aN.Analyzer
	if x == nil {
		fmt.Println("x is nil")
	}
}
