## consistentimports: detect inconsistent import aliases

[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/consistentimports)](https://goreportcard.com/report/github.com/nikolaydubina/consistentimports)

Report import paths and aliases count when same import path has multiple aliases. 

```bash
go install github.com/nikolaydubina/consistentimports@latest
```

```bash
consistentimports ./...
```

In example bellow same same import has two alises `passInspect` `!=` `passinspect`. This linter detects that.
```go
// apple.go
import (
    passInspect "golang.org/x/tools/go/analysis/passes/inspect"
)

// pear.go
import (
    passinspect "golang.org/x/tools/go/analysis/passes/inspect"
)
```

Example
```
-: "k8s.io/utils/net" netutils:4 netutil:1
-: "k8s.io/client-go/listers/core/v1" corelisters:1 listersv1:1 v1listers:1
-: "k8s.io/client-go/informers/core/v1" coreinformers:1 informers:1
-: "k8s.io/api/rbac/v1" rbacv1:4 v1:2
-: "k8s.io/apimachinery/pkg/runtime" runtime:3 kruntime:1
-: "k8s.io/api/imagepolicy/v1alpha1" imagepolicyv1alpha1:1 v1alpha1:1
```

## Alternatives

* https://github.com/julz/importas provides only fixed whitelist list or regex 
* https://staticcheck.io/docs/checks/ has no such linter

## Appendix A: packages in same module

As of `2023-06-08`, there seem to be no way to in `golang.org/x/tools/go/analysis` to:

* get list of packages in module
* check if two packages are in same module
* what is the current module

However, for practical use we really need to narrow down to current module, or else linter will keep checking recursively way outside of current module.

One simple heuristic is to check that package and the package it imports match long enough prefix.

## Appendix B: integration to linter aggregators

`staticcheck`
- `2023-06-08` proposal raised: https://github.com/dominikh/go-tools/issues/1413
