package proxy

import (
	"github.com/spf13/cobra"
)

func runStatusCommand(cmd *cobra.Command, _ []string) {
}

func initStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Returns the status of the reverse proxy.",
		Run:   runStatusCommand,
	}
}
