package tokens

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runListCommand(_ *cobra.Command, _ []string) error {
	tokens := All()
	if len(tokens) == 0 {
		return fmt.Errorf("no tokens created")
	}

	for _, token := range tokens {
		fmt.Println(strings.TrimSpace(string(token)))
	}

	return nil
}

func ListTokensCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:   "ls",
		Short: "list tokens",
	}, runListCommand, nil)
}
