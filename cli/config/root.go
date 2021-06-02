package config

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	proxyCommand := &cobra.Command{
		Use:   "config",
		Short: "Manage the config",
	}

	proxyCommand.AddCommand(initRunCommand())
	proxyCommand.AddCommand(initDeleteCommand())
	proxyCommand.AddCommand(initOpenCommand())

	return proxyCommand
}
