package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez/globals"
	"github.com/redwebcreation/hez/internal"
	"github.com/redwebcreation/hez/internal/ui"
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

var skipHealthchecks bool

func RunApplyCommand(_ *cobra.Command, _ []string) error {
	applications := internal.Config.Applications

	//var pool sync.WaitGroup
	//pool.Add(len(applications))

	//fatalErrors := make(chan error)
	//done := make(chan bool)

	for host, application := range applications {
		err := pullLatestImage(application)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		_, _ = application.StopContainer(internal.TemporaryContainer)

		_, err = application.CreateContainer(internal.TemporaryContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		err = ExecuteContainerDeployedHooks(application, internal.TemporaryContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		fmt.Printf("%s: container %s created\n", host, application.Name(internal.TemporaryContainer))

		err = WaitForContainerToBeHealthy(application, internal.TemporaryContainer)

		if err != nil {
			return err
		}

		_, _ = application.StopContainer(internal.ApplicationContainer)

		_, err = application.CreateContainer(internal.ApplicationContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		fmt.Printf("%s: new container %s created\n", host, application.Name(internal.ApplicationContainer))

		err = ExecuteContainerDeployedHooks(application, internal.ApplicationContainer)

		if err != nil {
			return err
		}

		err = WaitForContainerToBeHealthy(application, internal.ApplicationContainer)

		if err != nil {
			//fatalErrors <- err
			return err
		}

		_, err = application.StopContainer(internal.TemporaryContainer)

		fmt.Printf("%s: stopped temporary container %s\n", host, application.Name(internal.TemporaryContainer))

		if err != nil {
			//fatalErrors <- err
			return err
		}

		if *application.Warm {
			err := WarmContainer(application, internal.ApplicationContainer)

			if err != nil {
				return err
				//fatalErrors <- err
				//return
			}

			fmt.Printf("%s: container warmed up\n", host)
		}

		ui.Success(host + ": application is live")
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

func ExecuteContainerDeployedHooks(application *internal.Application, containerType int) error {
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

func WarmContainer(application *internal.Application, containerType int) error {
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
func WaitForContainerToBeHealthy(application *internal.Application, containerType int) error {
	if skipHealthchecks {
		return nil
	}

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
	return internal.CreateCommand(&cobra.Command{
		Use:   "apply",
		Short: "Applies your configuration to the server",
		Long:  `Applies your configuration to the server`,
	}, func(command *cobra.Command) {
		command.Flags().BoolVar(&skipHealthchecks, "skip-healthchecks", false, "Skip new container's healthchecks")
	}, RunApplyCommand)
}

func pullLatestImage(application *internal.Application) error {
	options := types.ImagePullOptions{}

	if !application.HasRegistry() {
		return nil
	}

	encodedAuth, _ := json.Marshal(map[string]string{
		"username": application.Registry.Username,
		"password": application.Registry.Password,
	})

	options.RegistryAuth = base64.StdEncoding.EncodeToString(encodedAuth)

	events, err := globals.Docker.ImagePull(context.Background(), application.Image, options)

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
