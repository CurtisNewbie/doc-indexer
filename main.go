package main

import (
	"os"

	"github.com/curtisnewbie/doc-indexer/docindexer"
)

func main() {
	docindexer.ServerRun(os.Args)
}
