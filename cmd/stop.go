package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/common"
	"github.com/wormable/nest/ansi")

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

		fmt.Printf("  %s- %s: %d containers stopped.%s\n", ansi.Gray.Fg(), application.Host, count, ansi.Reset)

		total += count
	}

	fmt.Printf("\n  %sStopped %d containers.%s\n", ansi.Green.Fg(), total, ansi.Reset)

	return nil
}

func StopCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "stop",
		Short: "Stops all containers.",
	}, runStopCommand, nil)
}
