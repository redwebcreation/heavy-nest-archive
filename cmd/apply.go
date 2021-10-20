package cmd

import (
	"fmt"

	"github.com/redwebcreation/nest/client"
	"github.com/redwebcreation/nest/internal"
	"github.com/redwebcreation/nest/ui"
	"github.com/spf13/cobra"
)

var skipHealthchecks bool

func runApplyCommand(_ *cobra.Command, args []string) error {
	if len(client.Config.Backends) == 0 {
		return fmt.Errorf("no backends configured")
	}

	if len(client.Config.Applications) == 0 {
		return fmt.Errorf("no applications found")
	}

	for _, application := range client.Config.Applications {
		if len(args) > 0 && application.Host != args[0] {
			fmt.Printf("- skipping %s\n", application.Host)
			continue
		}

		fmt.Printf("- %s\n", application.Host)

		application.SecondaryContainer().Deploy(client.DeploymentOptions{
			Pull:         true,
			Healthchecks: !skipHealthchecks,
		})
		application.PrimaryContainer().Deploy(client.DeploymentOptions{
			Pull:         false,
			Healthchecks: !skipHealthchecks,
		})

		application.SecondaryContainer().StopContainer()

		fmt.Printf("  %s%s deployed!%s\n", ui.Green.Fg(), application.Host, ui.Stop)
	}

	return nil
}

func ApplyCommand() *cobra.Command {
	return internal.CreateCommand(&cobra.Command{
		Use:   "apply [host]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Syncs the servers' state with your configuration",
	}, func(c *cobra.Command) {
		c.Flags().BoolVarP(&skipHealthchecks, "skip-healthchecks", "K", false, "Skip healthchecks")

	}, runApplyCommand)
}