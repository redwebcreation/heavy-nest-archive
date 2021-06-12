package cli

import (
	"github.com/redwebcreation/hez2/cli/proxy"
	"github.com/redwebcreation/hez2/globals"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	for _, arg := range os.Args {
		if arg == "--version" || arg == "-v" {
			globals.Ansi.Print("Hez " + globals.Version)
			return
		}
	}

	cli := &cobra.Command{
		Use:   "hez",
		Short: "Hez makes orchestrating containers easy.",
		Long:  `Hez is a tool to orchestrate containers and manage the environment around it.`,
	}

	cli.PersistentFlags().BoolP("version", "v", false, "Prints Hez's version.")

	cli.AddCommand(proxy.RootCommand())
	cli.AddCommand(ApplyCommand())
	cli.AddCommand(SelfUpdateCommand())
	cli.AddCommand(InfoCommand())
	cli.AddCommand(VersionCommand())

	cli.SilenceErrors = true

	globals.Ansi.Check(cli.Execute())
}
