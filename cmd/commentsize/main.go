// Command commentsize runs the comment-block-size analyzer standalone.
//
// As a go vet tool:
//
//	go build -o commentsize ./cmd/commentsize
//	go vet -vettool=$(pwd)/commentsize ./...
//	go vet -vettool=$(pwd)/commentsize -commentsize.max-lines=8 ./...
//
// Or directly:
//
//	go run ./cmd/commentsize -max-lines=8 ./...
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/smazurov/golangci-lint-commentsize"
)

func main() {
	singlechecker.Main(commentsize.New(commentsize.DefaultMaxLines))
}
