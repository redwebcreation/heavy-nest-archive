package certificates

import "github.com/spf13/cobra"

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "certificates",
	}

	cmd.AddCommand(InitCommand())

	return cmd
}