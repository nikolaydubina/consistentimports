package a

import (
	"fmt"
	str2 "strings"

	an2 "golang.org/x/tools/go/analysis"
)

func AA() {
	fmt.Println(str2.Join([]string{"a", "b", "c"}, ","))

	var x *an2.Analyzer
	if x == nil {
		fmt.Println("x is nil")
	}
}
