package core

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"os"
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
	env, err := FindEnvironment(application.Environment)

	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(env.Contents()), "\n")

	var envVariables []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		envVariables = append(envVariables, "-e", line)
	}

	for _, binding := range application.Bindings {
		envVariablesForBinding := envVariables

		envVariablesForBinding = append(envVariablesForBinding,
			"-e", "VIRTUAL_HOST="+binding.Host,
			"-e", "VIRTUAL_PORT="+binding.Port,
		)

		resp, err := GetDockerClient().ContainerCreate(context.Background(), &container.Config{
			Env:   envVariablesForBinding,
			Image: application.Image,
		}, &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}, nil, nil, getApplicationNameForBinding(application.Name, binding))

		if err != nil {
			return err
		}

		err = GetDockerClient().ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})

		if err != nil {
			return err
		}

		fmt.Println("  - Starting container [" + application.Image + "]" + " for " + binding.Host)
	}

	return nil
}

func (application Application) GetEnvironment() (Environment, error) {
	return FindEnvironment(application.Environment)
}

func (application Application) HasRunningContainers() bool {
	for _, applicationContainer := range application.GetContainers() {
		if applicationContainer.ID != "" {
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
	for _, applicationContainer := range application.GetContainers() {
		err := GetDockerClient().ContainerStop(context.Background(), applicationContainer.ID, nil)

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

	for _, applicationContainer := range application.GetContainers() {
		err := GetDockerClient().ContainerStop(context.Background(), applicationContainer.ID, nil)

		if err != nil {
			return err
		}

		fmt.Println("  - Stopping container [" + applicationContainer.ID + "]")

		err = GetDockerClient().ContainerRemove(context.Background(), applicationContainer.ID, types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		})

		if err != nil {
			return err
		}

		fmt.Println("  - Removed container [" + applicationContainer.ID + "]")
	}

	return nil
}
