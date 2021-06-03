package env

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	proxyCommand := &cobra.Command{
		Use:   "env",
		Short: "Manage the environments",
	}

	proxyCommand.AddCommand(initCommitCommand())
	proxyCommand.AddCommand(initListCommand())
	proxyCommand.AddCommand(initSyncCommand())

	return proxyCommand
}
