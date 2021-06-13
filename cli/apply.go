package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez2/ansi"
	"github.com/redwebcreation/hez2/core"
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

		fmt.Printf("%s: Temporary container created.\n", application.Host)

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

		fmt.Printf("%s: Container created.\n", application.Host)

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

		ansi.Text(application.Host+": Application is live.", ansi.Green)
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
					ansi.Loader(application.Host+": Warming up "+strconv.Itoa(i+1)+"/10", &counter)
					_, err = http.Get(proxiableContainer.Ipv4 + ":" + proxiableContainer.VirtualPort)
				}

				counter = 0
				ansi.Loader(application.Host+": Container warmed up.", &counter)
				break
			}
		}
	}

	return nil
}
func WaitForContainerToBeHealthy(containerId string, application globals.Application) {
	inspection, _ := inspectContainer(containerId)
	var counter int
	var secondsWaited int

	for isContainerStarting(inspection) {
		ansi.Loader(application.Host+": Waiting for container to be healthy ("+strconv.Itoa(secondsWaited)+"/"+strconv.FormatFloat(inspection.Config.Healthcheck.Interval.Seconds(), 'f', 0, 64)+")", &counter)
		inspection, _ = inspectContainer(containerId)
		time.Sleep(1 * time.Second)
		secondsWaited += 1
	}

	counter = 0
	ansi.Loader(application.Host+": Container is healthy.", &counter)
}

func inspectContainer(containerId string) (types.ContainerJSON, error) {
	inspection, err := globals.Docker.ContainerInspect(context.Background(), containerId)

	if err != nil {
		return inspection, err
	}

	return inspection, nil
}

func isContainerStarting(container types.ContainerJSON) bool {
	if container.State.Health == nil {
		return false
	}

	return container.State.Health.Status == "starting"
}

func ApplyCommand() *cobra.Command {
	return util.CreateCommand(&cobra.Command{
		Use:   "apply",
		Short: "Applies your configuration to the server.",
		Long:  `Applies your configuration to the server.`,
	}, nil, RunApplyCommand)
}

func pullLatestImage(application globals.Application) error {
	options := types.ImagePullOptions{}

	if application.HasRegistry() {
		encodedAuth, _ := json.Marshal(map[string]string{
			"username": application.Registry.Username,
			"password": application.Registry.Password,
		})

		options.RegistryAuth = base64.StdEncoding.EncodeToString(encodedAuth)
	}

	events, err := globals.Docker.ImagePull(context.Background(), application.Image, options)

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

		if first {
			fmt.Print("\n")
			first = false
		}

		newStatus := event.Status

		if status != newStatus {
			counter = 0
		}

		ansi.Loader(application.Host+": "+strings.Replace(event.Status, "Status: ", "", 1), &counter)

		status = newStatus
	}

	return nil
}
