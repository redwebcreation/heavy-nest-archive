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

func (application Application) Start(isTemporary bool) error {
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
		}, nil, nil, getApplicationNameForBinding(isTemporary, binding))

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

func getApplicationNameForBinding(isTemporary bool, binding Binding) string {
	name := strings.ReplaceAll(binding.Host, ".", "_") + "_" + binding.Port

	if isTemporary {
		name += "_temporary"
	}

	return name
}

type ShouldCleanup func(container types.Container) bool

func (application Application) CleanUpAllContainers() error {
	return application.CleanUp(func(container types.Container) bool {
		return true
	})
}

func (application Application) CleanUp(shouldCleanup ShouldCleanup) error {
	if !application.HasRunningContainers() {
		return nil
	}

	for _, applicationContainer := range application.GetContainers() {
		if !shouldCleanup(applicationContainer) {
			continue
		}

		err := GetDockerClient().ContainerStop(context.Background(), applicationContainer.ID, nil)

		if err != nil {
			return err
		}

		fmt.Println("  - Stopping container [" + applicationContainer.Names[0] + "]")

		err = GetDockerClient().ContainerRemove(context.Background(), applicationContainer.Names[0], types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		})

		if err != nil {
			return err
		}

		fmt.Println("  - Removed container [" + applicationContainer.Names[0] + "]")
	}

	return nil
}
