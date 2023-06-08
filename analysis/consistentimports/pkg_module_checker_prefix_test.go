package consistentimports_test

import (
	"fmt"
	"testing"

	"github.com/nikolaydubina/consistentimports/analysis/consistentimports"
)

func TestPrefixPkgModuleChecker_Same(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		num  uint
		same bool
	}{
		{
			a:    "github.com/nikolaydubina/consistentimports",
			b:    "github.com/nikolaydubina/consistentimports",
			num:  3,
			same: true,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a",
			b:    "github.com/nikolaydubina/consistentimports",
			num:  3,
			same: true,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "github.com/nikolaydubina/consistentimports/d/e/f/g",
			num:  3,
			same: true,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "uber.go/nikolaydubina/consistentimports/a/b/c/d",
			num:  3,
			same: false,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "github.com/nikolaydubina",
			num:  3,
			same: false,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "github.com/",
			num:  3,
			same: false,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "",
			num:  3,
			same: false,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "//",
			num:  3,
			same: false,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "/",
			num:  3,
			same: false,
		},
		{
			a:    "github.com/nikolaydubina/consistentimports/a/b/c/d",
			b:    "",
			num:  3,
			same: false,
		},
		{
			a:    "",
			b:    "",
			num:  3,
			same: false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d:(%s)(%s)", i, tt.a, tt.b), func(t *testing.T) {
			checker := consistentimports.PrefixPkgModuleChecker{NumSegments: tt.num}
			if got := checker.IsSameModule(tt.a, tt.b); got != tt.same {
				t.Errorf("PrefixPkgModuleChecker.IsSameModule() = %v, want %v", got, tt.same)
			}
		})
	}
}
