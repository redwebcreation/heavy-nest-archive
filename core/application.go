package core

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	Image         string `yaml:"image"`
	Host          string `yaml:"host"`
	ContainerPort int    `yaml:"container_port",omitempty`
	Env           []string
}

func (application Application) HasApplicationContainer() bool {
	container, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg(
			"name",
			application.Name(false)),
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

	env = append(env, "VIRTUAL_HOST="+application.Host, "VIRTUAL_PORT="+strconv.Itoa(application.ContainerPort))
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

func (application Application) Name(isEphemeral bool) string {
	name := strings.ReplaceAll(application.Host, ".", "_") + "_" + strconv.Itoa(application.ContainerPort)

	if isEphemeral {
		name += "_ephemeral"
	}

	return name
}

func (application Application) RemoveApplicationContainer() (string, error) {
	return RemoveContainer(application.Name(false))
}

func (application Application) RemoveEphemeralContainer() (string, error) {
	return RemoveContainer(application.Name(true))
}
