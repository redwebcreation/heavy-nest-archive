package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/common"
	"github.com/wormable/nest/ansi")

var skipHealthchecks bool

func runApplyCommand(_ *cobra.Command, args []string) error {
	if len(common.Config.Applications) == 0 {
		return fmt.Errorf("no applications found")
	}

	i := 0
	for _, application := range common.Config.Applications {
		if len(args) > 0 && application.Host != args[0] {
			fmt.Printf("- skipping %s\n", application.Host)
			i++
			continue
		}

		if i != 0 {
			fmt.Println()
		}
		fmt.Printf("  Deploying %s.\n", ansi.Blue.Fg()+application.Host+ansi.Reset)

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

		fmt.Printf("  %s is live!%s\n", ansi.Green.Fg()+application.Host, ansi.Reset)
		i++
	}

	return nil
}

func ApplyCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "apply [host]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Syncs the servers' state with your configuration",
	}, runApplyCommand, func(c *cobra.Command) {
		c.Flags().BoolVarP(&skipHealthchecks, "skip-healthchecks", "K", false, "Skip healthchecks")
	})
}
