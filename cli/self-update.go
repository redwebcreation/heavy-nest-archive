package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var draft bool
var prerelease bool

func RunSelfUpdateCommand(_ *cobra.Command, _ []string) error {
	fmt.Printf("Wants draft: %t\n", draft)
	fmt.Printf("Wants prerelease: %t\n", prerelease)
	return nil
}

func SelfUpdateCommand() *cobra.Command {
	command := CreateCommand(&cobra.Command{
		Use:   "self-update [version]",
		Short: "Updates Hez to the latest version.",
		Long:  `Updates Hez to the latest version or the one given as the first argument.`,
	}, func(command *cobra.Command) {
		command.Flags().BoolVarP(&draft, "edge", "e", false, "Updates to the main branch build.")
		command.Flags().BoolVarP(&prerelease, "prerelease", "p", false, "Updates to the latest prerelease.")
	}, RunSelfUpdateCommand)

	return command
}
