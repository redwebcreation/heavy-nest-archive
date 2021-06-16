package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez/ansi"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os/exec"
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
	applications := core.Config.Applications

	//var pool sync.WaitGroup
	//pool.Add(len(applications))

	//fatalErrors := make(chan error)
	//done := make(chan bool)

	for _, application := range applications {
		err := pullLatestImage(application)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		_, _ = application.StopTemporaryContainer()

		temporaryContainer, err := application.CreateTemporaryContainer()

		if err != nil {
			//fatalErrors <- err
			return err
		}

		fmt.Printf("%s: container %s created\n", application.Host, application.NameWithSuffix("temporary"))

		WaitForContainerToBeHealthy(temporaryContainer, application)

		_, err = application.StopApplicationContainer()

		if err != nil {
			//fatalErrors <- err
			return err
		}

		container, err := application.CreateApplicationContainer()

		if err != nil {
			//fatalErrors <- err
			return err
		}

		fmt.Printf("%s: new container %s created\n", application.Host, application.Name())

		WaitForContainerToBeHealthy(container, application)

		err = ExecuteContainerDeployedHooks(container, application)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		_, err = application.StopTemporaryContainer()

		fmt.Printf("%s: stopped temporary container %s\n", application.Host, application.NameWithSuffix("temporary"))

		if err != nil {
			//fatalErrors <- err
			return err
		}

		if *application.Warm {
			err := WarmContainer(application.Name())

			if err != nil {
				return err
				//fatalErrors <- err
				//return
			}

			fmt.Printf("%s: container warmed up\n", application.Host)
		}

		ansi.Text(application.Host+": application is live", ansi.Green)
		//pool.Done()
	}

	//go func() {
	//	pool.Wait()
	//	close(done)
	//}()
	//
	//select {
	//case <-done:
	//	break
	//case err := <-fatalErrors:
	//	close(fatalErrors)
	//	return err
	//}

	err := core.RefreshLastChangedTimestamp()

	if err != nil {
		return err
	}

	return nil
}

func ExecuteContainerDeployedHooks(container string, application core.Application) error {
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

func WarmContainer(container string) error {
	containers, err := core.GetProxiableContainers()
	if err != nil {
		return err
	}

	for _, proxiableContainer := range containers {
		if proxiableContainer.Container.ID == container {
			for i := 0; i < 10; i++ {
				_, err = http.Get(proxiableContainer.Ipv4 + ":" + proxiableContainer.VirtualPort)
			}

			break
		}
	}

	return nil
}
func WaitForContainerToBeHealthy(containerId string, application core.Application) {
	inspection, _ := inspectContainer(containerId)

	for isContainerStarting(inspection) {
		inspection, _ = inspectContainer(containerId)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("%s: %s is healthy.\n", application.Host, inspection.Name[1:])
}

func inspectContainer(containerId string) (types.ContainerJSON, error) {
	inspection, err := core.Docker.ContainerInspect(context.Background(), containerId)

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
	return core.CreateCommand(&cobra.Command{
		Use:   "apply",
		Short: "Applies your configuration to the server",
		Long:  `Applies your configuration to the server`,
	}, nil, RunApplyCommand)
}

func pullLatestImage(application core.Application) error {
	options := types.ImagePullOptions{}

	if application.HasRegistry() {
		encodedAuth, _ := json.Marshal(map[string]string{
			"username": application.Registry.Username,
			"password": application.Registry.Password,
		})

		options.RegistryAuth = base64.StdEncoding.EncodeToString(encodedAuth)
	}

	events, err := core.Docker.ImagePull(context.Background(), application.Image, options)

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(events)

	var event *Event
	var status string

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if status == event.Status {
			continue
		}

		fmt.Printf("%s: %s\n", application.Host, strings.ToLower(strings.Replace(event.Status, "Status: ", "", 1)))
		status = event.Status
	}

	return nil
}
