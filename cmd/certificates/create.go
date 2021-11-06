package certificates

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
	"github.com/wormable/nest/common"
)

func runCreateCommand(_ *cobra.Command, args []string) error {
	backend := args[0]

	for _, ip := range common.Config.Backends {
		if backend != ip {
			continue
		}

		return nil
	}

	return fmt.Errorf("could not find backend [%s]", backend)
}

func CreateCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use:  "create [backend]",
		Args: cobra.ExactArgs(1),
	}, nil, runCreateCommand)
}
