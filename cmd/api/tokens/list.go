package tokens

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runListCommand(_ *cobra.Command, _ []string) error {
	fmt.Println("here")
	return nil
}

func ListTokensCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:   "ls",
		Short: "list tokens",
	}, runListCommand, nil)
}
