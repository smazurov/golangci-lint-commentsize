package commentsize_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/smazurov/golangci-lint-commentsize"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), commentsize.New(3), "a")
}

// diagnosticsFor runs the analyzer over src and returns the raw diagnostics so
// a test can inspect the reported range, not just the message text.
func diagnosticsFor(t *testing.T, src string, maxLines int) ([]analysis.Diagnostic, *token.FileSet) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "x.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	var diags []analysis.Diagnostic
	az := commentsize.New(maxLines)
	pass := &analysis.Pass{
		Analyzer: az,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) { diags = append(diags, d) },
		ResultOf: map[*analysis.Analyzer]any{},
	}
	if _, err := az.Run(pass); err != nil {
		t.Fatal(err)
	}
	return diags, fset
}

// The flagged // block sits on lines 4-7.
const narrationSrc = `package b

func f() {
	// first line of narration
	// second line of narration
	// third line of narration
	// fourth line of narration
	_ = 0
}
`

func TestBareCountAnchoredAtStart(t *testing.T) {
	diags, fset := diagnosticsFor(t, narrationSrc, 3)
	if len(diags) != 1 {
		t.Fatalf("got %d diagnostics, want 1: %v", len(diags), diags)
	}
	d := diags[0]
	if strings.Contains(d.Message, "\n") {
		t.Errorf("message must stay single-line, got %q", d.Message)
	}
	if want := "comment block is 4 lines (max 3)"; d.Message != want {
		t.Errorf("message = %q, want %q", d.Message, want)
	}
	if line := fset.Position(d.Pos).Line; line != 4 {
		t.Errorf("anchored at line %d, want 4 (first line of the block)", line)
	}
}

// TestBlockCommentCounted covers a /* */ block: one ast.Comment spanning four
// physical lines, counted and reported at its first line.
func TestBlockCommentCounted(t *testing.T) {
	const src = `package b

func tall() {
	/* line one
	   line two
	   line three
	   line four */
	_ = 0
}
`
	diags, fset := diagnosticsFor(t, src, 3)
	if len(diags) != 1 {
		t.Fatalf("got %d diagnostics, want 1: %v", len(diags), diags)
	}
	if want := "comment block is 4 lines (max 3)"; diags[0].Message != want {
		t.Errorf("message = %q, want %q", diags[0].Message, want)
	}
	if line := fset.Position(diags[0].Pos).Line; line != 4 {
		t.Errorf("anchored at line %d, want 4", line)
	}
}
