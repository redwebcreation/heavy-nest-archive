package cli

import (
	"fmt"
	"github.com/redwebcreation/hez/ansi"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
)

func runStopCommand(_ *cobra.Command, _ []string) error {
	for _, application := range core.Config.Applications {
		current, _ := application.StopApplicationContainer()
		temporary, _ := application.StopTemporaryContainer()

		if current.ID != "" {
			fmt.Printf("Stopped %s (%s).\n", application.Name(), current.ID)
		}

		if temporary.ID != "" {
			fmt.Printf("Stopped %s (%s).\n", application.NameWithSuffix("temporary"), temporary.ID)
		}
	}
	ansi.Text("All container have been stopped successfully.", ansi.Green)

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
