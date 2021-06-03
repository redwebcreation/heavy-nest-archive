package env

import (
	"github.com/spf13/cobra"
)

func runListCommand(_ *cobra.Command, _ []string) {
}

func initListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all your configuration.",
		Run:   runListCommand,
	}
}
