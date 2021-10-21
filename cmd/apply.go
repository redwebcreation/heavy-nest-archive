package cmd

import (
	"fmt"

	"github.com/redwebcreation/nest/cmd/ui"
	"github.com/redwebcreation/nest/common"
	"github.com/spf13/cobra"
)

var skipHealthchecks bool

func runApplyCommand(_ *cobra.Command, args []string) error {
	if len(common.Config.Backends) == 0 {
		return fmt.Errorf("no backends configured")
	}

	if len(common.Config.Applications) == 0 {
		return fmt.Errorf("no applications found")
	}

	for _, application := range common.Config.Applications {
		if len(args) > 0 && application.Host != args[0] {
			fmt.Printf("- skipping %s\n", application.Host)
			continue
		}

		fmt.Printf("- %s\n", application.Host)

		application.Deploy(common.DeploymentOptions{
			Pull:         true,
			Healthchecks: !skipHealthchecks,
			Name:         application.TemporaryContainerName(),
		})
		application.Deploy(common.DeploymentOptions{
			Pull:         false,
			Healthchecks: !skipHealthchecks,
			Name:         application.ContainerName(),
		})

		application.StopContainer(application.TemporaryContainerName())

		fmt.Printf("  %s%s deployed!%s\n", ui.Green.Fg(), application.Host, ui.Stop)
	}

	return nil
}

func ApplyCommand() *cobra.Command {
	return CreateCommand(&cobra.Command{
		Use:   "apply [host]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Syncs the servers' state with your configuration",
	}, func(c *cobra.Command) {
		c.Flags().BoolVarP(&skipHealthchecks, "skip-healthchecks", "K", false, "Skip healthchecks")

	}, runApplyCommand)
}
