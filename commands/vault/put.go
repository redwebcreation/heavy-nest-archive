package vault

import (
	"github.com/spf13/cobra"
)

var force bool

func runPutCommand(_ *cobra.Command, args []string) error {
	return nil
}

func PutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put [key] [--force]",
		Args:  cobra.ExactArgs(1),
		Short: "Put a secret into the vault",
		RunE:  runPutCommand,
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite")

	return cmd
}
