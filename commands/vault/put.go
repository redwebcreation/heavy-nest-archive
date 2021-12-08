package vault

import (
	"fmt"
	"github.com/spf13/cobra"
)

func runPutCommand(_ *cobra.Command, args []string) error {
	_ = args[0]
	password, err := AskPrivately("Password: ")
	if err != nil {
		return err
	}

	value, err := AskPrivately("Value: ")
	if err != nil {
		return err
	}

	vault := Vault{
		Path: "/home/me/Code/server/dep/vault",
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
	return &cobra.Command{
		Use:   "put [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Put a secret into the vault",
		RunE:  runPutCommand,
	}
}
