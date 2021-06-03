package core

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"os"
	"strings"
)

type Binding struct {
	Host string `yaml:"host"`
	Port string `yaml:"container_port"`
}

type Application struct {
	Image    string `yaml:"image"`
	Bindings []Binding
	Env      string `yaml:"environment"`
}

func (application Application) Start() error {
	for _, binding := range application.Bindings {
		bindings := []string{"VIRTUAL_HOST=" + binding.Host, "VIRTUAL_PORT=" + binding.Port}

		resp, err := GetDockerClient().ContainerCreate(context.Background(), &container.Config{
			Env:   bindings,
			Image: application.Image,
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: EnvironmentPath(application.Env + "/current/.env"),
					Target: "/var/www/html/.env",
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}, nil, nil, getApplicationNameForBinding(binding))

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

func (application Application) GetStagingEnvironment() []byte {
	file, _ := os.ReadFile(EnvironmentPath(application.Env + "/staging/.env"))

	return file
}

func (application Application) GetCurrentEnvironment() []byte {
	file, _ := os.ReadFile(EnvironmentPath(application.Env + "/current/.env"))

	return file
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

func getApplicationNameForBinding(binding Binding) string {
	return strings.ReplaceAll(binding.Host, ".", "_") + "_" + binding.Port
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
