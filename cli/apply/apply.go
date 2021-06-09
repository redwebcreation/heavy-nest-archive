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

var skipHealthchecks bool
var withPulls bool

func run(_ *cobra.Command, _ []string) {
	config, _ := core.FindConfig(core.ConfigFile()).Resolve()

	fmt.Println("Estimated update time : " + GetEstimatedUpdateTime(config, skipHealthchecks) + "s")

	// Some space
	fmt.Println()

	var wg sync.WaitGroup
	wg.Add(len(config.Applications))

	for i := 0; i < len(config.Applications); i++ {
		go func(application core.Application) {
			if withPulls {
				PrintFor(application, "Pulling latest image.")
				err := application.PullLatestImage()

				if err != nil {
					fmt.Println(err)
				}
			}

			if !application.HasApplicationContainer() {
				container := application.CreateContainer(false)
				PrintFor(application, "Created container from ["+application.Image+"].")

				if !skipHealthchecks {
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
				}

				SuccessFor(application, "Application is live!")
				wg.Done()
				return
			}

			container := application.CreateContainer(true)

			PrintFor(application, "Created an ephemeral container from ["+application.Image+"]")

			if !skipHealthchecks {
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
			}

			application.RemoveApplicationContainer()

			PrintFor(application, "Stopped old container.")

			application.CreateContainer(false)

			PrintFor(application, "Created new container from ["+application.Image+"]")

			if !skipHealthchecks {
				for ContainerIsStarting(container) {
					time.Sleep(1 * time.Second)
				}

				if IsUnhealthy(container) {
					ErrorFor(application, "Container is is an unhealthy state.")
					PrintFor(application, "Can not roll back to the last healthy state.")
					wg.Done()
					return
				} else {
					PrintFor(application, "Container is in an healthy state.")
				}
			}

			application.RemoveEphemeralContainer()

			PrintFor(application, "Removed ephemeral container.")

			SuccessFor(application, "Application is live!")
			wg.Done()
		}(config.Applications[i])
	}

	wg.Wait()
}

func GetEstimatedUpdateTime(config core.ConfigData, skipHealthchecks bool) string {
	estimate := 0
	biggestInterval := 0

	for _, application := range config.Applications {
		inspection, _, err := core.GetDockerClient().ImageInspectWithRaw(context.Background(), application.Image)

		if err != nil {
			continue
		}

		if inspection.ContainerConfig.Healthcheck == nil {
			estimate += 500
		} else if !skipHealthchecks {
			interval := int(inspection.ContainerConfig.Healthcheck.Interval.Milliseconds())

			if interval > biggestInterval {
				biggestInterval = interval
			}
		}
	}

	estimate += biggestInterval

	duration, _ := time.ParseDuration(strconv.Itoa(estimate) + "ms")

	if duration.Seconds() <= 1 {
		duration, _ = time.ParseDuration("1s")
	}

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

	applyCmd.Flags().BoolVar(&skipHealthchecks, "skip-healthchecks", false, "Do no wait for containers to be healthy")
	applyCmd.Flags().BoolVar(&withPulls, "with-pulls", false, "Pull images to get the latest version")
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
