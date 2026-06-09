// Package a is a fixture. This package doc comment is long enough to exceed
// the test budget of three lines, yet it must NOT be flagged because it
// documents the package clause and is therefore exempt from the rule.
package a

// ExportedFn has a doc comment that runs well past the three-line budget on
// purpose. Doc comments on declarations are exempt, so no diagnostic should
// fire here even though this block is clearly taller than the limit allows.
func ExportedFn() {
	// one short why line is fine
	_ = 1

	// want `comment block is 4 lines`
	// this is a narration novel
	// that nobody asked for
	// and keeps going on
	_ = 2
}

// Config is exempt via its type doc comment, no matter how tall it grows here
// across several lines of prose that exceed the configured budget entirely.
type Config struct {
	// Field doc comments hang off ast.Field.Doc and are exempt as well even
	// when they stretch across more lines than the budget would allow here.
	Field int
}

func helper() {
	_ = 3 // a trailing comment stays on one line and is never flagged
}
