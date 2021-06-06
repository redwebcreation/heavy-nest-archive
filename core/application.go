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

type Application struct {
	Image         string `yaml:"image"`
	Host          string `yaml:"host"`
	ContainerPort string `yaml:"container_port"`
	Env           []string
}

func (application Application) HasApplicationContainer() bool {
	container, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", application.Name(false)),
		),
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return len(container) > 0
}

func (application Application) CreateContainer(isEphemeral bool) string {
	var env []string

	env = append(env, "VIRTUAL_HOST="+application.Host, "VIRTUAL_PORT="+application.ContainerPort)
	env = append(env, application.Env...)

	resp, err := GetDockerClient().ContainerCreate(context.Background(), &container.Config{
		Env:   env,
		Image: application.Image,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
	}, nil, nil, application.Name(isEphemeral))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = GetDockerClient().ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return resp.ID
}

//
//func (application Application) Start(binding Binding, isEphemeral bool) (string, error) {
//	var env []string
//
//	env = append(env, "VIRTUAL_HOST="+binding.Host, "VIRTUAL_PORT="+binding.ContainerPort)
//	env = append(env, application.Env...)
//
//	resp, err := GetDockerClient().ContainerCreate(context.Background(), &container.Config{
//		Env:   env,
//		Image: application.Image,
//	}, &container.HostConfig{
//		RestartPolicy: container.RestartPolicy{
//			Name: "unless-stopped",
//		},
//	}, nil, nil, binding.Name(isEphemeral))
//
//	if err != nil {
//		return "", err
//	}
//
//	err = GetDockerClient().ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
//
//	if err != nil {
//		return "", err
//	}
//
//	fmt.Println("  - Starting container [/" + binding.Name(isEphemeral) + "]")
//
//	return resp.ID, nil
//}
//
//func (application Application) HasRunningContainers() bool {
//	for _, applicationContainer := range application.GetContainers() {
//		if applicationContainer.ID != "" {
//			return true
//		}
//	}
//
//	return false
//}
//
//func (application Application) GetContainers() []types.Container {
//	containers, _ := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{
//		Filters: filters.NewArgs(filters.KeyValuePair{
//			Key:   "ancestor",
//			Value: application.Image,
//		}),
//	})
//
//	return containers
//}
//
//type ShouldCleanup func(container types.Container) bool
//
//func (application Application) CleanUpAllContainers() error {
//	return application.CleanUp(func(container types.Container) bool {
//		return true
//	})
//}
//
//func (application Application) CleanUp(shouldCleanup ShouldCleanup) error {
//	if !application.HasRunningContainers() {
//		return nil
//	}
//
//	for _, applicationContainer := range application.GetContainers() {
//		if !shouldCleanup(applicationContainer) {
//			continue
//		}
//
//		err := GetDockerClient().ContainerStop(context.Background(), applicationContainer.ID, nil)
//
//		if err != nil {
//			return err
//		}
//
//		fmt.Println("  - Stopping container [" + applicationContainer.Names[0] + "]")
//
//		err = GetDockerClient().ContainerRemove(context.Background(), applicationContainer.Names[0], types.ContainerRemoveOptions{
//			RemoveVolumes: false,
//			RemoveLinks:   false,
//			Force:         false,
//		})
//
//		if err != nil {
//			return err
//		}
//
//		fmt.Println("  - Removed container [" + applicationContainer.Names[0] + "]")
//	}
//
//	return nil
//}

func (application Application) Name(isEphemeral bool) string {
	name := strings.ReplaceAll(application.Host, ".", "_") + "_" + application.ContainerPort

	if isEphemeral {
		name += "_ephemeral"
	}

	return name
}

func (application Application) RemoveApplicationContainer() string {
	return RemoveContainer(application.Name(false))
}

func (application Application) RemoveEphemeralContainer() string {
	return RemoveContainer(application.Name(true))
}
