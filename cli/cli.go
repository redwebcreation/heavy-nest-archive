package cli

import (
	"fmt"
	"github.com/redwebcreation/hez2/ansi"
	"github.com/redwebcreation/hez2/globals"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	for _, arg := range os.Args {
		if arg == "--version" || arg == "-v" {
			ansi.Print("Hez " + globals.Version)
			return
		}
	}

	cli := &cobra.Command{
		Use:   "hez",
		Short: "Hez makes orchestrating containers easy.",
		Long:  `Hez is a tool to orchestrate containers and manage the environment around it.`,
	}

	cli.PersistentFlags().BoolP("version", "v", false, "Prints Hez's version.")

	cli.AddCommand(SelfUpdateCommand())
	cli.AddCommand(InfoCommand())
	cli.AddCommand(VersionCommand())

	ansi.Check(cli.Execute())
}

type CommandConfigurationHandler func(command *cobra.Command)
type CommandHandler func(_ *cobra.Command, _ []string) error

func CreateCommand(command *cobra.Command, commandConfigurationHandler CommandConfigurationHandler, Handler CommandHandler) *cobra.Command {
	if commandConfigurationHandler != nil {
		commandConfigurationHandler(command)
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")

		if showVersion {
			fmt.Println("Hez " + globals.Version)
			return nil
		}

		if Handler == nil {
			return nil
		}

		return Handler(cmd, args)
	}
	command.SilenceErrors = true
	command.SilenceUsage = true

	return command
}
