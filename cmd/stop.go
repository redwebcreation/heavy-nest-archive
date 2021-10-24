package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wormable/ui"
	"github.com/wormable/nest/common"
)

func runStopCommand(_ *cobra.Command, _ []string) error {
	fmt.Println()
	total := 0
	for _, application := range common.Config.Applications {
		stoppedPrimary := application.StopContainer(application.ContainerName())
		stoppedSecondary := application.StopContainer(application.TemporaryContainerName())

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
	return CreateCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stops all containers.",
	}, nil, runStopCommand)
}
