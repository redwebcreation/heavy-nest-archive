package apply

import (
	"context"
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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

	config, _ := configFile.Resolve()

	// Some space
	fmt.Println()

	fmt.Println("Estimated update time : " + GetEstimatedUpdateTime(config) + "s")

	// Some space
	fmt.Println()

	var wg sync.WaitGroup
	wg.Add(len(config.Applications))

	for i := 0; i < len(config.Applications); i++ {
		go func(application core.Application) {
			if !application.HasApplicationContainer() {
				container := application.CreateContainer(false)
				PrintFor(application, "Created container from ["+application.Image+"].")

				for ContainerIsStarting(container) {
					time.Sleep(1 * time.Second)
				}

				if IsUnhealthy(container) {
					ErrorFor(application, "Container is is an unhealthy state.")
					application.RemoveEphemeralContainer()
					ErrorFor(application, "Not rolling back as there's no healthy running container for this application.")
					wg.Done()
					return
				} else {
					PrintFor(application, "Container is in an healthy state.")
				}

				SuccessFor(application, "Application is live!")
				wg.Done()
				return
			}

			container := application.CreateContainer(true)

			PrintFor(application, "Created an ephemeral container from ["+application.Image+"]")

			for ContainerIsStarting(container) {
				time.Sleep(1 * time.Second)
			}

			if IsUnhealthy(container) {
				ErrorFor(application, "Container is is an unhealthy state.")
				application.RemoveEphemeralContainer()
				PrintFor(application, "Rolling back to the last healthy state.")
				wg.Done()
				return
			} else {
				PrintFor(application, "Container is in an healthy state.")
			}

			application.RemoveApplicationContainer()

			PrintFor(application, "Stopped old container.")

			application.CreateContainer(false)

			PrintFor(application, "Created new container from ["+application.Image+"]")

			for ContainerIsStarting(container) {
				time.Sleep(1 * time.Second)
			}

			PrintFor(application, "Container is in an healthy state.")

			application.RemoveEphemeralContainer()

			PrintFor(application, "Removed ephemeral container.")

			SuccessFor(application, "Application is live!")
			wg.Done()
		}(config.Applications[i])
	}

	wg.Wait()

	_ = core.SetKey("last_deployment", time.Now().String())
	_ = core.SetKey("previous_checksum", currentChecksum)
}

func GetEstimatedUpdateTime(config core.ConfigData) string {
	estimate := 0

	for _, application := range config.Applications {
		inspection, _, err := core.GetDockerClient().ImageInspectWithRaw(context.Background(), application.Image)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if inspection.ContainerConfig.Healthcheck == nil {
			estimate += 500
		} else {
			estimate += int(inspection.ContainerConfig.Healthcheck.Interval.Milliseconds())
		}
	}

	duration, _ := time.ParseDuration(strconv.Itoa(estimate) + "ms")

	return strconv.FormatFloat(duration.Seconds(), 'f', 0, 64)
}

func IsUnhealthy(container string) bool {
	inspection, err := core.GetDockerClient().ContainerInspect(context.Background(), container)

	if err != nil {
		fmt.Println(err)
		return true
	}

	if inspection.State.Health == nil {
		return false
	}

	return inspection.State.Health.Status == "unhealthy"
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

	return inspection.State.Health.Status != "healthy"
}

func PrintFor(application core.Application, message string) {
	fmt.Println(application.Host + ": " + message)
}

func ErrorFor(application core.Application, message string) {
	fmt.Println("\033[31m" + application.Host + ": " + message + "\033[0m")
}

func SuccessFor(application core.Application, message string) {
	fmt.Println("\033[32m" + application.Host + ": " + message + "\033[0m")
}
