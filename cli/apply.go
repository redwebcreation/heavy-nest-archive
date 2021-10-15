package cli

import (
	"fmt"

	"github.com/redwebcreation/hez/client"
	"github.com/redwebcreation/hez/internal"
	"github.com/spf13/cobra"
)

var skipHealthchecks bool

func runApplyCommand(_ *cobra.Command, _ []string) error {
	if len(client.Config.Backends) == 0 {
		return fmt.Errorf("no backends configured")
	}

	if len(client.Config.Applications) == 0 {
		return fmt.Errorf("no applications found")
	}

	for host, application := range client.Config.Applications {

	}

	return nil
}

func ApplyCommand() *cobra.Command {
	return internal.CreateCommand(&cobra.Command{
		Use:   "apply",
		Short: "Syncs the servers' state with your configuration",
	}, func(c *cobra.Command) {
		c.Flags().BoolVarP(&skipHealthchecks, "skip-healthchecks", "K", false, "Skip healthchecks")

	}, runApplyCommand)
}
