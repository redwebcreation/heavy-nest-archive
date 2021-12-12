package main

import (
	"fmt"
	"os"

	"github.com/redwebcreation/nest/commands/vault"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/globals"
	"github.com/spf13/cobra"
)

func main() {
	remote := config.Remote{
		Url: "git@github.com:redwebcreation/server.git",
	}

	os.Exit(0)
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

	err = cli.Execute()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
