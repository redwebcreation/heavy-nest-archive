package tokens

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runRevokeCommand(_ *cobra.Command, args []string) error {
	for _, rawToken := range args {
		token := strings.TrimSpace(rawToken)

		if token == "" {
			continue
		}

		Revoke(token)
		fmt.Printf("%s\n", token)
	}

	return nil
}

func RevokeTokenCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:   "revoke [token] [...token]",
		Short: "revoke a new token",
		Args:  cobra.MinimumNArgs(1),
	}, runRevokeCommand, nil)
}
