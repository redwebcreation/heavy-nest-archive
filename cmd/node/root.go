package node

import (
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "WILL BE MOVED TO ITS OWN CLI",
	}

	cmd.AddCommand(ListenCommand())

	return cmd
}
