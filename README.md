# golangci-lint-commentsize

A [go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) linter that
flags contiguous **non-doc comment blocks** taller than a configurable line
budget — the "comment novel" smell. Doc comments attached to a declaration
(function, type, var, const, the package clause, struct/interface fields) are
exempt, since Go conventions ask for those and other linters already govern
them. The target is free-floating narration: banner blocks and multi-line
prose inside function bodies.

```go
func process() {
	// This block explains, at length, what the following ten lines do,
	// step by step, restating each statement in prose, which is exactly
	// the kind of narration this linter exists to discourage, because the
	// code below already says all of this on its own.   <-- flagged
	...
}
```

Both `//` line-comment runs and multi-line `/* */` blocks are measured by their
source-line span.

## Use as a golangci-lint module plugin (recommended)

golangci-lint cannot load an external analyzer from config alone — it must be
compiled in. The [module plugin system](https://golangci-lint.run/docs/plugins/module-plugins/)
builds a drop-in `custom-gcl` binary that pulls this module remotely.

`.custom-gcl.yml`:

```yaml
version: v2.12.2          # match your installed golangci-lint
name: custom-gcl
destination: ./bin
plugins:
  - module: github.com/smazurov/golangci-lint-commentsize
    version: v0.1.0
```

`.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - commentsize
  settings:
    custom:
      commentsize:
        type: module
        description: Flags contiguous non-doc comment blocks taller than the line budget
        settings:
          max-lines: 6
```

Then:

```bash
golangci-lint custom      # builds ./bin/custom-gcl with this plugin baked in
./bin/custom-gcl run
```

## Use standalone (go vet)

```bash
go install github.com/smazurov/golangci-lint-commentsize/cmd/commentsize@latest
go vet -vettool=$(which commentsize) ./...
# override the budget:
commentsize -max-lines=8 ./...
```

## Use as a library

```go
import "github.com/smazurov/golangci-lint-commentsize"

az := commentsize.New(6) // *analysis.Analyzer
```

## Settings

| Setting     | Default | Meaning                                                        |
|-------------|---------|----------------------------------------------------------------|
| `max-lines` | `6`     | Maximum source-line span of a non-doc comment block before it is flagged. |

## License

MIT — see [LICENSE](LICENSE).
