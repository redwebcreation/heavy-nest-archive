package tokens

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runRevokeCommand(_ *cobra.Command, _ []string) error {
	fmt.Println("here")
	return nil
}

func RevokeTokenCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:   "revoke [name]",
		Short: "revoke a new token",
	}, runRevokeCommand, nil)
}
