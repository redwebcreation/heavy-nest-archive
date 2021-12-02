package tokens

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runCreateCommand(_ *cobra.Command, _ []string) error {
	token := Create()

	fmt.Print(*token)

	return nil
}

func CreateTokenCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:   "create",
		Short: "Create a new token",
	}, runCreateCommand, nil)
}
