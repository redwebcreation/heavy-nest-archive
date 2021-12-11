package vault

import (
	"github.com/spf13/cobra"
)

func runGetCommand(_ *cobra.Command, args []string) error {
	return nil
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Get secrets from the vault",
		RunE:  runGetCommand,
	}

	return cmd
}
