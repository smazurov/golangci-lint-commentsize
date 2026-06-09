// Package commentsize provides a go/analysis pass that flags contiguous
// comment blocks taller than a configurable line budget.
//
// Doc comments attached to a declaration (function, type, var, const, the
// package clause) are exempt: Go conventions ask for those and other linters
// already govern them. The target is free-floating narration — banner blocks
// and multi-line "novels" inside function bodies.
package commentsize

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// DefaultMaxLines is the line budget used when New is given maxLines <= 0.
const DefaultMaxLines = 6

// New returns an analyzer that flags any non-doc comment block taller than
// maxLines source lines. maxLines <= 0 falls back to DefaultMaxLines.
//
// The diagnostic anchors at the comment's first line with a single-line,
// bare-count message. It deliberately does not echo the comment text or set an
// end position: golangci-lint discards analysis.Diagnostic.End (its buildIssues
// keeps only the start Pos), so a range never reaches an annotation, and a
// multi-line message would break the GitHub ::error:: command. The annotation
// already sits on the comment, so the count is enough.
//
// maxLines is also exposed as a -max-lines flag for standalone go vet use.
func New(maxLines int) *analysis.Analyzer {
	if maxLines <= 0 {
		maxLines = DefaultMaxLines
	}
	a := &checker{maxLines: maxLines}
	az := &analysis.Analyzer{
		Name: "commentsize",
		Doc:  "flags contiguous non-doc comment blocks taller than the configured line budget",
		Run:  a.run,
	}
	az.Flags.IntVar(&a.maxLines, "max-lines", maxLines,
		"maximum line height of a non-doc comment block before it is flagged")
	return az
}

type checker struct {
	maxLines int
}

func (c *checker) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		docs := docComments(file)
		for _, cg := range file.Comments {
			if docs[cg] {
				continue
			}
			start := pass.Fset.Position(cg.Pos()).Line
			end := pass.Fset.Position(cg.End()).Line
			if lines := end - start + 1; lines > c.maxLines {
				pass.Reportf(cg.Pos(), "comment block is %d lines (max %d)", lines, c.maxLines)
			}
		}
	}
	return nil, nil
}

// docComments returns the set of CommentGroups that document a declaration.
// Each *ast.CommentGroup in File.Comments is the same pointer stored on a
// node's Doc field, so identity comparison exempts them reliably. Every AST
// node carrying a Doc field is visited, including struct/interface fields.
func docComments(file *ast.File) map[*ast.CommentGroup]bool {
	docs := make(map[*ast.CommentGroup]bool)
	add := func(cg *ast.CommentGroup) {
		if cg != nil {
			docs[cg] = true
		}
	}
	add(file.Doc)
	ast.Inspect(file, func(n ast.Node) bool {
		switch d := n.(type) {
		case *ast.FuncDecl:
			add(d.Doc)
		case *ast.GenDecl:
			add(d.Doc)
		case *ast.TypeSpec:
			add(d.Doc)
		case *ast.ValueSpec:
			add(d.Doc)
		case *ast.ImportSpec:
			add(d.Doc)
		case *ast.Field:
			add(d.Doc)
		}
		return true
	})
	return docs
}
