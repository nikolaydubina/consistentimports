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

var Analyzer = &analysis.Analyzer{
	Name:      "consistentimports",
	Doc:       "detect inconsistent import aliases, reports import paths and aliases count",
	Run:       run,
	Flags:     flag.FlagSet{},
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{new(pathAliasCountsFact)},
}

type pathAliasCountsFact map[[2]string]uint

func (f *pathAliasCountsFact) AFact() {}

func run(pass *analysis.Pass) (interface{}, error) {
	var pathAliasCount = pathAliasCountsFact{}

	for _, file := range pass.Files {
		for _, spec := range file.Imports {
			// upstream packages
			var fact pathAliasCountsFact
			if pass.ImportPackageFact(imported(pass.TypesInfo, spec), &fact) {
				for k, v := range fact {
					pathAliasCount[k] += v
				}
			}

			if spec == nil {
				continue
			}

			var alias string
			if spec.Name != nil {
				alias = spec.Name.Name
			}

			if alias == "" {
				continue
			}

			// path already wrapped into ""
			path := spec.Path.Value

			pathAliasCount[[2]string{path, alias}]++
		}
	}

	var pathAliasAgg map[string][]aliasCount = map[string][]aliasCount{}
	for k, v := range pathAliasCount {
		pathAliasAgg[k[0]] = append(pathAliasAgg[k[0]], aliasCount{Alias: k[1], Count: v})
	}
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

	// sort paths by aliases count
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

	return nil, nil
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
