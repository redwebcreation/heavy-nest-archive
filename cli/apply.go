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

		_, _ = application.StopContainer(core.TemporaryContainer)

		_, err = application.CreateContainer(core.TemporaryContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		err = ExecuteContainerDeployedHooks(application, core.TemporaryContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		fmt.Printf("%s: container %s created\n", application.Host, application.Name(core.TemporaryContainer))

		err = WaitForContainerToBeHealthy(application, core.TemporaryContainer)

		if err != nil {
			return err
		}

		_, _ = application.StopContainer(core.ApplicationContainer)
		_, err = application.CreateContainer(core.ApplicationContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		fmt.Printf("%s: new container %s created\n", application.Host, application.Name(core.ApplicationContainer))

		err = ExecuteContainerDeployedHooks(application, core.ApplicationContainer)

		if err != nil {
			return err
		}

		err = WaitForContainerToBeHealthy(application, core.ApplicationContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		_, err = application.StopContainer(core.TemporaryContainer)

		fmt.Printf("%s: stopped temporary container %s\n", application.Host, application.Name(core.TemporaryContainer))

		if err != nil {
			//fatalErrors <- err
			return err
		}

		if *application.Warm {
			err := WarmContainer(application, core.ApplicationContainer)

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

	return nil
}

func ExecuteContainerDeployedHooks(application core.Application, containerType int) error {
	if len(application.Hooks.ContainerDeployed) == 0 {
		return nil
	}

	container, err := application.GetContainer(containerType)

	if err != nil {
		return err
	}

	for _, c := range application.Hooks.ContainerDeployed {
		command := []string{
			"exec",
			container.Ref.ID,
		}

		for _, piece := range strings.Split(c, " ") {
			command = append(command, piece)
		}

		cmd := exec.Command("docker", command...)

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		fmt.Println(application.Host + ": running command [" + strings.Join(command[2:], " ") + "]")

		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}

		if stdout.Len() > 0 && err != nil {
			return errors.New(stdout.String())
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func WarmContainer(application core.Application, containerType int) error {
	container, err := application.GetContainer(containerType)

	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		_, err = http.Get("http://" + container.Ip + ":" + container.Port)
		if err != nil {
			return err
		}
	}

	return nil
}
func WaitForContainerToBeHealthy(application core.Application, containerType int) error {
	container, err := application.GetContainer(containerType)

	if err != nil {
		return err
	}

	inspection, _ := inspectContainer(container.Ref.ID)

	i := -1.0
	for isContainerStarting(inspection) {
		inspection, _ = inspectContainer(container.Ref.ID)

		if i == -1 {
			i = inspection.Config.Healthcheck.Interval.Seconds()
		}

		fmt.Println("Waiting for container to be healthy (" + strconv.FormatFloat(i, 'f', 0, 64) + "/" + strconv.FormatFloat(inspection.Config.Healthcheck.Interval.Seconds(), 'f', 0, 64) + ")")
		time.Sleep(1 * time.Second)
		i--
	}

	inspection, _ = inspectContainer(container.Ref.ID)

	if inspection.State.Health == nil || inspection.State.Health.Status == "healthy" {
		fmt.Printf("%s: %s is healthy.\n", application.Host, inspection.Name[1:])
		return nil
	}

	_, _ = application.StopContainer(containerType)
	return errors.New("container is unhealthy, rolling back")
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
