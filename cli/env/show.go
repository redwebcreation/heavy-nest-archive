package env

import (
	"github.com/spf13/cobra"
)

var staging bool
var current bool

func runShowCommand(_ *cobra.Command, _ []string) {
}

func initShowCommand() *cobra.Command {
	showCommand := &cobra.Command{
		Use:   "show",
		Short: "Shows a specific env file",
		Run:   runShowCommand,
	}

	showCommand.Flags().BoolVar(&staging, "staging", false, "Only the staging environment")
	showCommand.Flags().BoolVar(&current, "current", false, "Only the current environment")

	return showCommand
}
