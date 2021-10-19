package cmd

import (
	"fmt"
	"strings"

	"github.com/redwebcreation/nest/client"
	"github.com/redwebcreation/nest/internal"
	"github.com/redwebcreation/nest/ui"
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

	for _, application := range client.Config.Applications {
		ui.Title("    " + application.Host)

		application.SecondaryContainer().Stop()
		application.SecondaryContainer().Start()

		application.PrimaryContainer().Stop()
		application.PrimaryContainer().Start()

		fmt.Println("\n    Application deployed!")
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

func getContainerBaseName(a client.Application, host string) string {
	return strings.ReplaceAll(host, ".", "_") + "_" + a.Port
}
