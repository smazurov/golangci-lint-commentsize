// Command commentsize runs the comment-block-size analyzer standalone.
//
// As a go vet tool:
//
//	go build -o commentsize ./cmd/commentsize
//	go vet -vettool=$(pwd)/commentsize ./...
//
// Or directly, overriding the budget:
//
//	go run ./cmd/commentsize -max-lines=3 ./...
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/smazurov/golangci-lint-commentsize"
)

func main() {
	singlechecker.Main(commentsize.New(commentsize.DefaultMaxLines))
}
