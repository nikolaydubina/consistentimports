package a

import (
	"fmt"
	str "strings"

	an "golang.org/x/tools/go/analysis"
)

//go:generate mockery --name DB
type DB interface {
	Get(key string) (an.Analyzer, error)
}

func A() {
	fmt.Println(str.Join([]string{"a", "b", "c"}, ","))

	var x *an.Analyzer
	if x == nil {
		fmt.Println("x is nil")
	}
}
