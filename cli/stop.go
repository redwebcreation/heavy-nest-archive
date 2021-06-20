package cli

import (
	"fmt"
	"github.com/redwebcreation/hez/ansi"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
)

func runStopCommand(_ *cobra.Command, _ []string) error {
	for host, application := range core.Config.Applications {
		current, _ := application.StopContainer(core.ApplicationContainer)
		temporary, _ := application.StopContainer(core.TemporaryContainer)

		if current.Ref != nil {
			fmt.Printf("Stopped %s (%s).\n", application.Name(core.ApplicationContainer), host)
		}

		if temporary.Ref != nil {
			fmt.Printf("Stopped %s (%s).\n", application.Name(core.TemporaryContainer), host)
		}
	}
	ansi.Success("All container have been stopped successfully.")

	return nil
}

func StopCommand() *cobra.Command {
	return core.CreateCommand(
		&cobra.Command{
			Use:   "stop",
			Short: "Stops all running containers.",
			Long:  `Stops all running containers.`,
		}, nil, runStopCommand)
}
