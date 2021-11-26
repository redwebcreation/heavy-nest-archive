package tokens

import "github.com/spf13/cobra"

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage API tokens",
	}

	cmd.AddCommand(
		CreateTokenCommand(),
		ListTokensCommand(),
		RevokeTokenCommand(),
	)

	return cmd
}
