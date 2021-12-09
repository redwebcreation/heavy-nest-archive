package vault

import (
	"fmt"
	"github.com/spf13/cobra"
)

var force bool

func runPutCommand(_ *cobra.Command, args []string) error {
	vault := Vault{
		Path: "/home/me/Code/server/dep/vault",
	}

	if vault.Has(args[0]) && !force {
		return fmt.Errorf("secret already exists")
	}

	password, err := AskPrivately("Password: ")
	if err != nil {
		return err
	}

	value, err := AskPrivately("Value: ")
	if err != nil {
		return err
	}

	err = vault.Put(Secret{
		Key:      args[0],
		Value:    value,
		Password: password,
	})
	if err != nil {
		return err
	}

	fmt.Println("OK")
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
