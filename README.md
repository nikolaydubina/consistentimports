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

Example
```
-: "k8s.io/utils/net" netutils:4 netutil:1
-: "k8s.io/client-go/listers/core/v1" corelisters:1 listersv1:1 v1listers:1
-: "k8s.io/client-go/informers/core/v1" coreinformers:1 informers:1
-: "k8s.io/api/rbac/v1" rbacv1:4 v1:2
-: "k8s.io/apimachinery/pkg/runtime" runtime:3 kruntime:1
-: "k8s.io/api/imagepolicy/v1alpha1" imagepolicyv1alpha1:1 v1alpha1:1
-: "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis
```

## Alternatives

* https://github.com/julz/importas provides only fixed whitelist list or regex 
* https://staticcheck.io/docs/checks/ has no such linter
