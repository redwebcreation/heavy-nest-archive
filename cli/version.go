package cli

import (
	"github.com/redwebcreation/hez2/core"
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
)

var short bool

func RunVersionCommand(_ *cobra.Command, _ []string) error {
	if short {
		core.Ansi.Print(globals.Version)
	} else {
		core.Ansi.Print("Hez " + globals.Version)
	}

	return nil
}

func VersionCommand() *cobra.Command {
	return util.CreateCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hez",
		Long:  `All software has versions. This is Hez's.`,
	}, func(command *cobra.Command) {
		command.Flags().BoolVarP(&short, "short", "s", false, "Prints out only the version number.")
	}, RunVersionCommand)
}
