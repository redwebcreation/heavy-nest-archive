package core

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"os"
	"strings"
)

type Application struct {
	Image         string   `yaml:"image"`
	Host          string   `yaml:"host"`
	ContainerPort string   `yaml:"container_port"`
	Network       string   `yaml:"network"`
	Warm          *bool    `yaml:"warm"`
	Env           []string `yaml:"env"`
	Volumes       []struct {
		From string
		To   string
	} `yaml:"volumes"`
	Registry struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
	}
	//Replicas      int
	Hooks struct {
		ContainerDeployed []string `yaml:"container_deployed"`
	} `yaml:"hooks"`
}

type Container struct {
	Type int
	Ip   string
	Port string
	Name string
	Ref  *types.Container
}

const (
	ApplicationContainer = iota
	TemporaryContainer
	AnyType
)

func (application Application) Name(containerType int) string {
	baseName := strings.ReplaceAll(application.Host, ".", "_") + "_" + application.ContainerPort

	if containerType == AnyType {
		application, _ := application.GetContainer(ApplicationContainer)

		if application.Ref.ID != "" {
			return baseName
		}

		return baseName + "_temporary"
	}

	if containerType == ApplicationContainer {
		return baseName
	}

	return baseName + "_temporary"
}

func (application Application) CreateContainer(containerType int) (string, error) {
	var mounts []mount.Mount
	env := []string{
		"VIRTUAL_HOST=" + application.Host,
		"VIRTUAL_PORT=" + application.ContainerPort,
	}

	networkDetails, _ := application.GetNetwork()

	for _, volume := range application.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: volume.From,
			Target: volume.To,
		})
	}

	resolvedEnvironment, err := ResolveEnvironmentVariables(env, application.Env)

	if err != nil {
		return "", err
	}

	resp, err := Docker.ContainerCreate(context.Background(), &container.Config{
		Env:   resolvedEnvironment,
		Image: application.Image,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		Mounts: mounts,
	}, nil, nil, application.Name(containerType))

	if err != nil {
		return "", err
	}

	err = Docker.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})

	if err != nil {
		return "", err
	}

	_ = Docker.NetworkDisconnect(context.Background(), networkDetails.ID, resp.ID, true)

	err = Docker.NetworkConnect(context.Background(), networkDetails.ID, resp.ID, nil)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (application Application) GetNetwork() (types.NetworkResource, error) {
	networks, err := Docker.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return types.NetworkResource{}, err
	}

	var network types.NetworkResource

	for _, maybeNetwork := range networks {
		if maybeNetwork.Name == application.Network {
			network = maybeNetwork
			break
		}
	}

	networkDetails, err := Docker.NetworkInspect(context.Background(), network.ID, types.NetworkInspectOptions{})

	if err != nil {
		return networkDetails, err
	}

	return networkDetails, nil
}

func (application Application) HasRegistry() bool {
	return application.Registry.Server != "" &&
		application.Registry.Username != "" &&
		application.Registry.Password != ""
}

func (application Application) GetContainer(containerType int) (Container, error) {
	containers, err := Docker.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
	}

	for _, c := range containers {
		name := application.Name(containerType)

		if "/"+name == c.Names[0] {
			return Container{
				Type: containerType,
				Ip:   c.NetworkSettings.Networks[application.Network].IPAddress,
				Port: application.ContainerPort,
				Name: name,
				Ref:  &c,
			}, err
		}

	}

	return Container{}, errors.New("no container found")
}

func (application Application) StopContainer(containerType int) (Container, error) {
	c, err := application.GetContainer(containerType)

	if err != nil {
		return Container{}, err
	}

	err = Docker.ContainerStop(context.Background(), c.Ref.ID, nil)

	if err != nil {
		return c, err
	}

	err = Docker.ContainerRemove(context.Background(), c.Name, types.ContainerRemoveOptions{})

	if err != nil {
		return c, err
	}

	return c, nil
}

func ResolveEnvironmentVariables(variables []string, env []string) ([]string, error) {
	for _, envVariable := range env {
		if strings.Contains(envVariable, "=") {
			variables = append(variables, envVariable)
		} else {
			contents, err := os.ReadFile(envVariable)
			if err != nil {
				return nil, err
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

	return variables, nil
}
