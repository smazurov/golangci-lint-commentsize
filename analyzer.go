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
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Default line budget for a non-doc comment block before it is flagged.
const DefaultMaxLines = 6

// New returns an analyzer that flags any non-doc comment block spanning more
// than maxLines source lines. A maxLines <= 0 falls back to DefaultMaxLines.
//
// The budget is also exposed as a -max-lines flag on the analyzer, so a
// standalone singlechecker/go vet build can override it without recompiling.
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
				pass.Reportf(cg.Pos(),
					"comment block is %d lines (max %d): say WHY in one line, don't narrate\n%s",
					lines, c.maxLines, commentText(cg))
			}
		}
	}
	return nil, nil
}

// commentText reproduces the block as written, keeping the // and /* */
// markers. cg.Text() is unsuitable here: it strips markers and reflows, so it
// would not echo the comment the author actually wrote.
func commentText(cg *ast.CommentGroup) string {
	var b strings.Builder
	for i, c := range cg.List {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(c.Text)
	}
	return b.String()
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
