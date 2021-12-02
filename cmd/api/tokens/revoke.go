package tokens

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runRevokeCommand(_ *cobra.Command, args []string) error {
	raw := args[0]
	tokens := strings.Split(raw, "\n")

	for _, rawToken := range tokens {
		token := strings.TrimSpace(rawToken)

		if token == "" {
			continue
		}

		Revoke(token)
		fmt.Printf("Revoked %s\n", token)
	}

	return nil
}

func RevokeTokenCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:   "revoke [name]",
		Short: "revoke a new token",
	}, runRevokeCommand, nil)
}
