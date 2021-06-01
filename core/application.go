package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"os"
	"os/exec"
	"strings"
)

type Binding struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Application struct {
	Name        string `yaml:"name"`
	Image       string `yaml:"image"`
	Bindings    []Binding
	Environment string `yaml:"environment"`
}

func (application Application) Start() error {
	environment, err := FindEnvironment(application.Environment)

	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(environment.Contents()), "\n")

	var environmentVariables []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		environmentVariables = append(environmentVariables, "-e", line)
	}

	for _, binding := range application.Bindings {
		dockerCommand := []string{"run", "-d", "--restart", "unless-stopped", "--name", getApplicationNameForBinding(application.Name, binding)}

		dockerCommand = append(dockerCommand, environmentVariables...)
		dockerCommand = append(dockerCommand,
			"-e", "VIRTUAL_HOST="+binding.Host,
			"-e", "VIRTUAL_PORT="+binding.Port,
			application.Image,
		)

		var stdErr bytes.Buffer
		cmd := exec.Command("docker", dockerCommand...)
		cmd.Stderr = &stdErr
		err = cmd.Run()

		if stdErr.Len() > 0 {
			return errors.New(stdErr.String())
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (application Application) GetEnvironment() (Environment, error) {
	return FindEnvironment(application.Environment)
}

func (application Application) HasRunningContainers() bool {
	for _, container := range application.GetContainers() {
		if container.ID != "" {
			return true
		}
	}

	return false
}

func (application Application) GetContainers() []types.Container {
	containers, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "ancestor",
			Value: application.Image,
		}),
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return containers
}

func (application Application) StopContainers() error {
	for _, container := range application.GetContainers() {
		err := GetDockerClient().ContainerStop(context.Background(), container.ID, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

func getApplicationNameForBinding(name string, binding Binding) string {
	return strings.ReplaceAll(binding.Host, ".", "_") + "_" + binding.Port + "_" + name
}

func (application Application) CleanUp() error {
	if !application.HasRunningContainers() {
		return nil
	}

	for _, container := range application.GetContainers() {
		err := GetDockerClient().ContainerStop(context.Background(), container.ID, nil)

		if err != nil {
			return err
		}

		fmt.Println("  - Stopping container [" + container.ID + "]")

		err = GetDockerClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		})

		if err != nil {
			return err
		}

		fmt.Println("  - Removed container [" + container.ID + "]")
	}

	return nil
}
