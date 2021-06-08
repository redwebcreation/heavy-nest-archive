package main

import (
	"fmt"
	"github.com/redwebcreation/hez/cli"
	box "github.com/redwebcreation/hez/core/embed"
	"os"
)

func main() {
	for _, arg := range os.Args {
		if arg == "--version" {
			fmt.Printf("Hez %s", string(box.Get("/version")))
			return
		}
	}

	cli.Execute()
}
