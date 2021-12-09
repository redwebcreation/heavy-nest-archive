package vault

import (
	"fmt"
	"github.com/spf13/cobra"
)

func runGetCommand(_ *cobra.Command, args []string) error {
	vault := Vault{
		Resolver: FileStorage{
			Path: "/home/me/Code/server/dep/vault",
		},
	}

	password, err := AskPrivately("Password: ")
	if err != nil {
		return err
	}

	secret := Secret{
		Key:      args[0],
		Password: password,
	}

	decrypted, err := vault.Get(&secret)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", decrypted)

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
