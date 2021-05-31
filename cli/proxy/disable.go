package proxy

import (
	"github.com/spf13/cobra"
)

func runDisableCommand(cmd *cobra.Command, _ []string) {

}

func initDisableCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disables the reverse proxy.",
		Run:   runDisableCommand,
	}
}
