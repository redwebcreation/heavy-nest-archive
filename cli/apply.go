package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez2/core"
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type Event struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func RunApplyCommand(_ *cobra.Command, _ []string) error {
	applications := globals.Config.Applications

	if len(applications) == 0 {
		return errors.New("no applications found in the configuration")
	}

	for _, application := range globals.Config.Applications {
		err := pullLatestImage(application)

		if err != nil {
			return err
		}

		_, _ = application.StopTemporaryContainer()

		temporaryContainer, err := application.CreateTemporaryContainer()

		core.Ansi.Printf("%s: Temporary container created.\n", application.Host)

		if err != nil {
			return err
		}

		WaitForContainerToBeHealthy(temporaryContainer, application)

		_, err = application.StopApplicationContainer()

		if err != nil {
			return err
		}

		container, err := application.CreateApplicationContainer()

		if err != nil {
			return err
		}

		core.Ansi.Printf("%s: Container created.\n", application.Host)

		WaitForContainerToBeHealthy(container, application)

		err = ExecuteContainerDeployedHooks(container, application)

		if err != nil {
			return err
		}

		_, err = application.StopTemporaryContainer()

		if err != nil {
			return err
		}

		if *application.Warm {
			err := WarmContainer(container, application)

			if err != nil {
				return err
			}
		}

		core.Ansi.Success(application.Host + ": Application is live.")
	}

	err := core.RefreshLastApplyExecution()

	if err != nil {
		return err
	}

	return nil
}

func ExecuteContainerDeployedHooks(container string, application globals.Application) error {
	if len(application.Hooks.ContainerDeployed) == 0 {
		return nil
	}

	// TODO: I couldn't make it work using the standard docker client.
	for _, c := range application.Hooks.ContainerDeployed {
		command := []string{
			"exec",
			container,
		}

		for _, piece := range strings.Split(c, " ") {
			command = append(command, piece)
		}

		cmd := exec.Command("docker", command...)

		var stderr bytes.Buffer

		cmd.Stderr = &stderr

		err := cmd.Run()

		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func WarmContainer(container string, application globals.Application) error {
	if core.IsProxyEnabled() {
		for i := 0; i < 10; i++ {
			_, err := http.Get(application.Host)

			if err != nil {
				return err
			}
		}
	} else {
		containers, err := core.GetProxiableContainers()
		if err != nil {
			return err
		}

		for _, proxiableContainer := range containers {
			if proxiableContainer.Container.ID == container {
				counter := 0
				for i := 0; i < 10; i++ {
					core.Ansi.StatusLoader(application.Host+": Warming up "+strconv.Itoa(i+1)+"/10", &counter)
					_, err = http.Get(proxiableContainer.Ipv4 + ":" + proxiableContainer.VirtualPort)
				}

				counter = 0
				core.Ansi.StatusLoader(application.Host+": Container warmed up.", &counter)
				break
			}
		}
	}

	return nil
}
func WaitForContainerToBeHealthy(containerId string, application globals.Application) {
	starting, _ := isContainerStarting(containerId)
	var counter int

	core.Ansi.NewLine()

	for starting {
		core.Ansi.StatusLoader(application.Host+": Waiting for container to be healthy", &counter)
		starting, _ = isContainerStarting(containerId)
	}

	counter = 0
	core.Ansi.StatusLoader(application.Host+": Container is healthy.", &counter)
}

func isContainerStarting(containerId string) (bool, error) {
	inspection, err := globals.Docker.ContainerInspect(context.Background(), containerId)

	if err != nil {
		return false, err
	}

	if inspection.State.Health == nil {
		return false, nil
	}

	return inspection.State.Health.Status != "healthy", nil
}

func ApplyCommand() *cobra.Command {
	return util.CreateCommand(&cobra.Command{
		Use:   "apply",
		Short: "Applies your configuration to the server.",
		Long:  `Applies your configuration to the server.`,
	}, nil, RunApplyCommand)
}

func pullLatestImage(application globals.Application) error {
	events, err := globals.Docker.ImagePull(context.Background(), application.Image, types.ImagePullOptions{
	})

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(events)

	var event *Event
	var first = true
	var counter int
	var status string

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if !core.Ansi {
			continue
		}

		if first {
			fmt.Print("\n")
			first = false
		}

		newStatus := event.Status

		if status != newStatus {
			counter = 0
		}

		core.Ansi.StatusLoader(application.Host+": "+strings.Replace(event.Status, "Status: ", "", 1), &counter)

		status = newStatus
	}

	if !core.Ansi {
		core.Ansi.Printf("%s: Pulled out the latest image for %s", application.Host, application.Image)
	}
	return nil
}
