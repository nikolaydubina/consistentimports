package consistentimports

import (
	"go/ast"
	"regexp"
)

const regexGeneratedComment = `^// Code generated .* DO NOT EDIT\.$`

var regexGeneratedCommentMatcher = regexp.MustCompile(regexGeneratedComment)

// isFileGenerated based on Go conventions
// https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source
func isFileGenerated(f *ast.File) bool {
	for _, q := range f.Comments {
		if q == nil {
			continue
		}
		for _, c := range q.List {
			if c == nil {
				continue
			}
			if regexGeneratedCommentMatcher.Match([]byte(c.Text)) {
				return true
			}
		}
	}
	return false
}
