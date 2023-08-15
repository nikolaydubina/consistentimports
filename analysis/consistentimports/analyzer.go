package consistentimports

import (
	"flag"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Analyzer for consistentimports.
// Inspired by analysis facts example: https://cs.opensource.google/go/x/tools/+/refs/tags/v0.9.3:go/analysis/passes/pkgfact/pkgfact.go
var Analyzer = &analysis.Analyzer{
	Name:      "consistentimports",
	Doc:       "detect inconsistent import aliases, reports import paths and aliases count",
	Run:       run,
	Flags:     flag.FlagSet{},
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{new(pathAliasCountsFact)},
}

var (
	skipGeneratedFiles bool
	numPrefixPathSame  uint
)

func init() {
	Analyzer.Flags.BoolVar(&skipGeneratedFiles, "skip-generated", true, `skip generated files`)
	Analyzer.Flags.UintVar(&numPrefixPathSame, "num-prefix-path-same", 3, `number of same prefix segments split by "/" in import path to consider it belonging to same module (e.g. "github.com/user/repo/pkga" matches "github.com/user/repo/pkgb")`)
}

type pathAliasCountsFact map[[2]string]uint

func newPathAliasCountsFact() pathAliasCountsFact { return pathAliasCountsFact{} }

func (f *pathAliasCountsFact) AFact() {}

func run(pass *analysis.Pass) (interface{}, error) {
	pathAliasCount := newPathAliasCountsFact()

	pkgModuleChecker := PrefixPkgModuleChecker{NumSegments: numPrefixPathSame}

	for _, file := range pass.Files {
		if ast.IsGenerated(file) && skipGeneratedFiles {
			continue
		}

		for _, spec := range file.Imports {
			// merge stats from upstream packages in same module
			if pkgModuleChecker.IsSameModule(pass.Pkg.Path(), spec.Path.Value) {
				if upfact := newPathAliasCountsFact(); pass.ImportPackageFact(imported(pass.TypesInfo, spec), &upfact) {
					for k, v := range upfact {
						pathAliasCount[k] += v
					}
				}
			}

			if spec == nil {
				continue
			}

			var alias string
			if spec.Name != nil {
				alias = spec.Name.Name
			}

			if alias == "" || alias == "_" {
				continue
			}

			// path already wrapped into ""
			path := spec.Path.Value

			pathAliasCount[[2]string{path, alias}]++
		}
	}

	reportPathAliasCount(pass, pathAliasCount)

	if len(pathAliasCount) > 0 {
		pass.ExportPackageFact(&pathAliasCount)
	}

	return nil, nil
}

func reportPathAliasCount(pass *analysis.Pass, pathAliasCount pathAliasCountsFact) {
	var pathAliasAgg map[string][]aliasCount = map[string][]aliasCount{}
	for k, v := range pathAliasCount {
		pathAliasAgg[k[0]] = append(pathAliasAgg[k[0]], aliasCount{Alias: k[1], Count: v})
	}

	// sort aliases by count for each path from highest to lowest
	for k := range pathAliasAgg {
		sort.Slice(pathAliasAgg[k], func(i, j int) bool {
			ci := pathAliasAgg[k][i].Count
			cj := pathAliasAgg[k][j].Count
			if ci == cj {
				return pathAliasAgg[k][i].Alias < pathAliasAgg[k][j].Alias
			}
			return ci > cj
		})
	}

	// sort paths by aliases count from highest to lowest
	pathsByCount := make([]string, 0, len(pathAliasAgg))
	for k := range pathAliasAgg {
		pathsByCount = append(pathsByCount, k)
	}
	sort.Slice(pathsByCount, func(i, j int) bool {
		li := len(pathAliasAgg[pathsByCount[i]])
		lj := len(pathAliasAgg[pathsByCount[j]])
		if li == lj {
			return pathsByCount[i] < pathsByCount[j]
		}
		return li > lj
	})

	// in each path report aliases
	for _, path := range pathsByCount {
		aliases := pathAliasAgg[path]
		if len(aliases) == 1 {
			continue
		}
		pass.Reportf(token.NoPos, `%s %s`, path, printPathAliasCount(aliases))
	}
}

type aliasCount struct {
	Alias string
	Count uint
}

func printPathAliasCount(vs []aliasCount) string {
	var b strings.Builder
	for i, v := range vs {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(v.Alias)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(int(v.Count)))
	}
	return b.String()
}

func imported(info *types.Info, spec *ast.ImportSpec) *types.Package {
	obj, ok := info.Implicits[spec]
	if !ok {
		obj = info.Defs[spec.Name] // renaming import
	}
	return obj.(*types.PkgName).Imported()
}
