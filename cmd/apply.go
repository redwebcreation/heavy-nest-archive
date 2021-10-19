package cmd

import (
	"fmt"
	"strings"

	"github.com/redwebcreation/nest/client"
	"github.com/redwebcreation/nest/internal"
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
		temporaryContainer := application.GetDeploymentConfigurationFor(
			getContainerBaseName(application, host) + "_temporary",
		)
		temporaryContainer.Deploy()

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