package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/commands/vault"
	"github.com/wormable/nest/globals"
	"os"
)

func main() {
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("nest %s\n", globals.Version)
			os.Exit(0)
		}
	}

	cli := cobra.Command{
		Use:   "nest",
		Short: "A simple yet powerful container orchestrator",
	}

	cli.AddCommand(vault.RootCommand())

	err := cli.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
