package tokens

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runRevokeAllCommand(_ *cobra.Command, args []string) error {
	for _, t := range All() {
		Revoke(string(t))
		fmt.Printf("%s\n", t)
	}

	return nil
}

func RevokeAllTokenCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:   "revoke-all",
		Short: "revoke all tokens",
	}, runRevokeAllCommand, nil)
}
