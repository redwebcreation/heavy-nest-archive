package core

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	Image         string `yaml:"image"`
	Host          string `yaml:"host"`
	ContainerPort int    `yaml:"container_port",omitempty`
	Env           []string
	Volumes       []struct {
		From string `yaml:"from"`
		To   string `yaml:"to"`
	}
}

func (application Application) HasApplicationContainer() bool {
	containers, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg(
			"name",
			application.Name(false)),
		),
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return len(containers) > 0
}

func (application Application) CreateContainer(isEphemeral bool) string {
	config, _ := FindConfig(ConfigFile()).Resolve()

	var mounts []mount.Mount
	env := []string{
		"VIRTUAL_HOST=" + application.Host,
		"VIRTUAL_PORT=" + strconv.Itoa(application.ContainerPort),
	}

	networkDetails, _ := FindNetwork(config.Network)

	for _, volume := range application.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: volume.From,
			Target: volume.To,
		})
	}

	fmt.Println(networkDetails.Name)

	resp, err := GetDockerClient().ContainerCreate(context.Background(), &container.Config{
		Env:   ResolveEnvironmentVariables(env, application.Env),
		Image: application.Image,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		Mounts: mounts,
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkDetails.ID: {
				NetworkID: networkDetails.ID,
			},
		},
	}, nil, application.Name(isEphemeral))

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

func ResolveEnvironmentVariables(variables []string, env []string) []string {
	for _, envVariable := range env {
		if strings.Contains(envVariable, "=") {
			variables = append(variables, envVariable)
		} else {
			contents, err := os.ReadFile(envVariable)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, envFileVariable := range strings.Split(string(contents), "\n") {
				trimmed := strings.TrimSpace(envFileVariable)

				if trimmed == "" {
					continue
				}

				variables = append(variables, trimmed)
			}
		}
	}

	return variables
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

func (application Application) PullLatestImage() error {
	_, err := GetDockerClient().ImagePull(context.Background(), application.Image, types.ImagePullOptions{})

	return err
}
