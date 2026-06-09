package commentsize_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/smazurov/golangci-lint-commentsize"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), commentsize.New(3), "a")
}

// TestBlockComment covers what analysistest cannot: a multi-line /* */ block,
// whose want directive would otherwise be parsed as the expectation itself.
func TestBlockComment(t *testing.T) {
	const src = `package b

func tall() {
	/* line one
	   line two
	   line three
	   line four */
	_ = 0
}

func short() {
	/* one liner */
	_ = 0
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "b.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	var msgs []string
	pass := &analysis.Pass{
		Analyzer: commentsize.New(3),
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) { msgs = append(msgs, d.Message) },
		ResultOf: map[*analysis.Analyzer]any{},
	}
	if _, err := commentsize.New(3).Run(pass); err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("got %d diagnostics, want 1: %v", len(msgs), msgs)
	}
	if want := "comment block is 4 lines (max 3)"; !contains(msgs[0], want) {
		t.Errorf("message %q does not contain %q", msgs[0], want)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
