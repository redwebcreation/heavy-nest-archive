package internal

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/redwebcreation/hez/globals"
	"os"
	"strings"
)

type Application struct {
	Image         string `yaml:"image"`
	Host          string
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

func (a Application) Name(containerType int) string {
	baseName := strings.ReplaceAll(a.Host, ".", "_") + "_" + a.ContainerPort

	if containerType == AnyType {
		application, _ := a.GetContainer(ApplicationContainer)

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

func (a Application) CreateContainer(containerType int) (string, error) {
	var mounts []mount.Mount
	env := []string{
		"VIRTUAL_HOST=" + a.Host,
		"VIRTUAL_PORT=" + a.ContainerPort,
	}

	networkDetails, err := a.GetNetwork()

	if err != nil {
		return "", err
	}

	for _, volume := range a.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: volume.From,
			Target: volume.To,
		})
	}

	resolvedEnvironment, err := ResolveEnvironmentVariables(env, a.Env)

	if err != nil {
		return "", err
	}

	resp, err := globals.Docker.ContainerCreate(context.Background(), &container.Config{
		Env:   resolvedEnvironment,
		Image: a.Image,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		Mounts: mounts,
	}, nil, nil, a.Name(containerType))

	if err != nil {
		return "", err
	}

	err = globals.Docker.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})

	if err != nil {
		return "", err
	}

	_ = globals.Docker.NetworkDisconnect(context.Background(), networkDetails.ID, resp.ID, true)

	err = globals.Docker.NetworkConnect(context.Background(), networkDetails.ID, resp.ID, nil)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (a Application) GetNetwork() (types.NetworkResource, error) {
	networks, err := globals.Docker.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return types.NetworkResource{}, err
	}

	var network types.NetworkResource

	for _, maybeNetwork := range networks {
		if maybeNetwork.Name == a.Network {
			network = maybeNetwork
			break
		}
	}

	networkDetails, err := globals.Docker.NetworkInspect(context.Background(), network.ID, types.NetworkInspectOptions{})

	if err != nil {
		return networkDetails, errors.New("could not find network " + a.Network + " (" + err.Error() + ")")
	}

	return networkDetails, nil
}

func (a Application) HasRegistry() bool {
	return a.Registry.Server != "" &&
		a.Registry.Username != "" &&
		a.Registry.Password != ""
}

func (a Application) GetContainer(containerType int) (Container, error) {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		return Container{}, err
	}

	for _, c := range containers {
		name := a.Name(containerType)

		if "/"+name == c.Names[0] {
			found := Container{
				Type: containerType,
				Port: a.ContainerPort,
				Name: name,
				Ref:  &c,
			}

			if c.NetworkSettings != nil && c.NetworkSettings.Networks[a.Network] != nil {
				found.Ip = c.NetworkSettings.Networks[a.Network].IPAddress
			}

			return found, err
		}

	}

	return Container{}, errors.New("no container found")
}

func (a Application) StopContainer(containerType int) (Container, error) {
	c, err := a.GetContainer(containerType)

	if err == nil {
		_ = globals.Docker.ContainerStop(context.Background(), c.Ref.ID, nil)
	}

	err = globals.Docker.ContainerRemove(context.Background(), a.Name(containerType), types.ContainerRemoveOptions{})

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
