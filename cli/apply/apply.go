package apply

import (
	"context"
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"sync"
	"time"
)

func run(cmd *cobra.Command, _ []string) {
	configFile := core.FindConfig(core.ConfigFile())

	currentChecksum, _ := configFile.Checksum()
	previousChecksum, _ := core.GetKey("previous_checksum")
	fmt.Println("Previous config checksum : " + previousChecksum)
	fmt.Println("Current config checksum : " + currentChecksum)

	force, _ := cmd.Flags().GetBool("force")

	if currentChecksum == previousChecksum {
		if force {
			fmt.Println("No changes. Not aborting as --force is set to true.")
		} else {
			fmt.Println("No changes. Aborting.")
			return
		}
	} else {
		fmt.Println("Found changes.")
	}

	// Some space
	fmt.Println()

	config, _ := configFile.Resolve()

	var wg sync.WaitGroup
	wg.Add(len(config.Applications))

	for i := 0; i < len(config.Applications); i++ {
		go func(application core.Application) {
			if !application.HasApplicationContainer() {
				container := application.CreateContainer(false)

				for ContainerIsStarting(container) {
					time.Sleep(1 * time.Second)
				}

				wg.Done()
				return
			}

			container := application.CreateContainer(true)

			for ContainerIsStarting(container) {
				time.Sleep(1 * time.Second)
			}

			application.RemoveApplicationContainer()

			application.CreateContainer(false)

			for ContainerIsStarting(container) {
				time.Sleep(1 * time.Second)
			}

			application.RemoveEphemeralContainer()

			wg.Done()
		}(config.Applications[i])
	}

	wg.Wait()
	
	_ = core.SetKey("previous_checksum", currentChecksum)
}

func NewCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Applies your configuration to the server.",
		Run:   run,
	}

	applyCmd.Flags().BoolP("force", "f", false, "Force the apply command to run")

	return applyCmd
}

func ContainerIsStarting(containerId string) bool {
	inspection, err := core.GetDockerClient().ContainerInspect(context.Background(), containerId)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if inspection.State.Health == nil {
		return false
	}

	return inspection.State.Health.Status == "starting"
}
