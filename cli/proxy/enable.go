package proxy

import (
	"github.com/spf13/cobra"
)

func runEnableCommand(cmd *cobra.Command, _ []string) {

}

func initEnableCommand() *cobra.Command {
	enableCommand := &cobra.Command{
		Use:   "enable",
		Short: "Enables the reverse proxy.",
		Run:   runEnableCommand,
	}

	return enableCommand
}
