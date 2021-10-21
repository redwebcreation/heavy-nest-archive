package cmd

import (
	"fmt"

	"github.com/redwebcreation/nest/client"
	"github.com/redwebcreation/nest/internal"
	"github.com/redwebcreation/nest/ui"
	"github.com/spf13/cobra"
)

func runStopCommand(_ *cobra.Command, _ []string) error {
	fmt.Println()
	total := 0
	for _, application := range client.Config.Applications {
		stoppedPrimary := application.PrimaryContainer().StopContainer()
		stoppedSecondary := application.SecondaryContainer().StopContainer()

		count := 0

		if stoppedPrimary {
			count++
		}

		if stoppedSecondary {
			count++
		}

		fmt.Printf("  %s- %s: %d containers stopped.%s\n", ui.Gray.Fg(), application.Host, count, ui.Stop)

		total += count
	}

	fmt.Printf("\n  %sStopped %d containers.%s\n", ui.Green.Fg(), total, ui.Stop)

	return nil
}

func StopCommand() *cobra.Command {
	return internal.CreateCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stops all containers.",
	}, nil, runStopCommand)
}
