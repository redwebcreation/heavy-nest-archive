package vault

import "github.com/spf13/cobra"

func runDeleteCommand(_ *cobra.Command, _ []string) error {
	return nil
}

func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a secret from the vault",
		RunE:  runDeleteCommand,
	}

	return cmd
}
