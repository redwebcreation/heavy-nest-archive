package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
	"io"
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

		globals.Ansi.Printf("%s: Temporary container created.\n", application.Host)

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

		globals.Ansi.Printf("%s: Container created.\n", application.Host)

		WaitForContainerToBeHealthy(container, application)

		_, err = application.StopTemporaryContainer()

		if err != nil {
			return err
		}

		globals.Ansi.Success(application.Host + ": Application is live.")
	}

	return nil
}

func WaitForContainerToBeHealthy(containerId string, application globals.Application) {
	starting, _ := isContainerStarting(containerId)
	var counter int

	globals.Ansi.NewLine()

	for starting {
		globals.Ansi.StatusLoader(application.Host+": Waiting for container to be healthy", &counter)
		starting, _ = isContainerStarting(containerId)
	}

	counter = 0
	globals.Ansi.StatusLoader(application.Host+": Container is healthy.", &counter)
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
	events, err := globals.Docker.ImagePull(context.Background(), application.Image, types.ImagePullOptions{})
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

		if !globals.Ansi {
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

		globals.Ansi.StatusLoader(application.Host+": "+strings.Replace(event.Status, "Status: ", "", 1), &counter)

		status = newStatus
	}

	if !globals.Ansi {
		globals.Ansi.Printf("%s: Pulled out the latest image for %s", application.Host, application.Image)
	}
	return nil
}
