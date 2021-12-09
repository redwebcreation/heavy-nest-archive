package vault

import "github.com/spf13/cobra"

func RootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "vault",
		Short: "manage your secrets",
	}

	root.AddCommand(
		PutCommand(),
		GetCommand(),
		//DeleteCommand(),
		//ListCommand(),
		//MetadataCommand(),
	)

	return root
}
