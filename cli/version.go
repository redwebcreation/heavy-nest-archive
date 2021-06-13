package cli

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/redwebcreation/hez/globals"
	"github.com/spf13/cobra"
)

var short bool

func RunVersionCommand(_ *cobra.Command, _ []string) error {
	if short {
		fmt.Println(globals.Version)
	} else {
		fmt.Println("Hez " + globals.Version)
	}

	return nil
}

func VersionCommand() *cobra.Command {
	return core.CreateCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hez",
		Long:  `All software has versions. This is Hez's.`,
	}, func(command *cobra.Command) {
		command.Flags().BoolVarP(&short, "short", "s", false, "Prints out only the version number.")
	}, RunVersionCommand)
}
