package a

import (
	"fmt"
	str "strings"

	an "golang.org/x/tools/go/analysis"
)

func A() {
	fmt.Println(str.Join([]string{"a", "b", "c"}, ","))

	var x *an.Analyzer
	if x == nil {
		fmt.Println("x is nil")
	}
}
