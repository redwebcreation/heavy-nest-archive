package proxy

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	proxyCommand := &cobra.Command{
		Use: "proxy",
	}

	proxyCommand.AddCommand(initEnableCommand())
	proxyCommand.AddCommand(initDisableCommand())
	proxyCommand.AddCommand(initStatusCommand())

	return proxyCommand
}
