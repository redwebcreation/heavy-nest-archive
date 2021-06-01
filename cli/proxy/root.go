package proxy

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	proxyCommand := &cobra.Command{
		Use:   "proxy",
		Short: "Manage the integrated lightning fast reverse proxy",
	}

	proxyCommand.AddCommand(initEnableCommand())
	proxyCommand.AddCommand(initDisableCommand())
	proxyCommand.AddCommand(initStatusCommand())
	proxyCommand.AddCommand(initRunCommand())

	return proxyCommand
}
