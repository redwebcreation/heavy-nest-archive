package proxy

import (
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cli := &cobra.Command{
		Use:   "proxy",
		Short: "Manage the reverse proxy.",
		Long:  `Manage the reverse proxy.`,
	}

	cli.AddCommand(RunCommand())
	cli.AddCommand(DisableCommand())
	cli.AddCommand(EnableCommand())
	cli.AddCommand(StatusCommand())

	return cli
}
