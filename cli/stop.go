package cli

import (
	"fmt"
	"github.com/redwebcreation/hez/internal"
	ui "github.com/redwebcreation/hez/internal/ui"
	"github.com/spf13/cobra"
)

func runStopCommand(_ *cobra.Command, _ []string) error {
	for host, application := range internal.Config.Applications {
		current, _ := application.StopContainer(internal.ApplicationContainer)
		temporary, _ := application.StopContainer(internal.TemporaryContainer)

		if current.Ref != nil {
			fmt.Printf("Stopped %s (%s).\n", application.Name(internal.ApplicationContainer), host)
		}

		if temporary.Ref != nil {
			fmt.Printf("Stopped %s (%s).\n", application.Name(internal.TemporaryContainer), host)
		}
	}
	ui.Success("All container have been stopped successfully.")

	return nil
}

func StopCommand() *cobra.Command {
	return internal.CreateCommand(
		&cobra.Command{
			Use:   "stop",
			Short: "Stops all running containers.",
			Long:  `Stops all running containers.`,
		}, nil, runStopCommand)
}
