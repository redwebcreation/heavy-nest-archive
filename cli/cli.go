package cli

import (
	"fmt"
	"github.com/redwebcreation/hez/ansi"
	"github.com/redwebcreation/hez/cli/proxy"
	"github.com/redwebcreation/hez/globals"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	for _, arg := range os.Args {
		if arg == "--version" || arg == "-v" {
			fmt.Println("Hez " + globals.Version)
			return
		}
	}

	cli := &cobra.Command{
		Use:   "hez",
		Short: "Hez makes orchestrating containers easy",
		Long:  `Hez is a tool to orchestrate containers and manage the environment around it`,
	}

	cli.PersistentFlags().BoolP("version", "v", false, "Prints Hez's version")

	cli.AddCommand(proxy.RootCommand())
	cli.AddCommand(ApplyCommand())
	cli.AddCommand(SelfUpdateCommand())
	cli.AddCommand(StopCommand())
	cli.AddCommand(InfoCommand())
	cli.AddCommand(VersionCommand())

	cli.SilenceErrors = true

	err := cli.Execute()

	if err != nil {
		ansi.Text(err.Error(), ansi.Red)
		os.Exit(1)
	}
}
