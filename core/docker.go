package core

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strings"
)

func GetDockerClient() *client.Client {
	cli, _ := client.NewClientWithOpts(client.FromEnv)

	return cli
}

type ProxiableContainer struct {
	Name        string
	Ipv4        string
	VirtualHost string
	VirtualPort string
	Container   *types.ContainerJSON
}

func GetWhitelistedDomains() []string {
	var domains []string

	proxiableContainers, _ := GetProxiableContainers()

	for _, proxiableContainer := range proxiableContainers {
		domains = append(domains, proxiableContainer.VirtualHost)
	}

	return domains
}

func GetProxiableContainers() ([]ProxiableContainer, error) {
	networks, err := GetDockerClient().NetworkList(context.Background(), types.NetworkListOptions{})
	var proxiableContainers []ProxiableContainer

	if err != nil {
		return nil, err
	}

	var networkBridge types.NetworkResource

	for _, network := range networks {
		if network.Name == "bridge" {
			networkBridge = network
		}
	}

	bridgeDetails, err := GetDockerClient().NetworkInspect(context.Background(), networkBridge.ID, types.NetworkInspectOptions{})

	if err != nil {
		return nil, err
	}

	for _, bridgeContainer := range bridgeDetails.Containers {
		containerList, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{})

		if err != nil {
			return nil, err
		}

		if len(containerList) != 1 {
			return nil, errors.New("zero or more than one container found")
		}

		container := containerList[0]
		inspectedContainer, err := GetDockerClient().ContainerInspect(context.Background(), container.ID)

		if err != nil {
			return nil, err
		}

		var virtualHost string
		virtualPort := "80"

		for _, envVariable := range inspectedContainer.Config.Env {
			if strings.HasPrefix(envVariable, "VIRTUAL_HOST=") {
				virtualHost = strings.SplitAfter(envVariable, "VIRTUAL_HOST=")[1]
			}

			if strings.HasPrefix(envVariable, "VIRTUAL_PORT=") {
				virtualPort = strings.SplitAfter(envVariable, "VIRTUAL_PORT=")[1]
			}
		}

		if virtualHost == "" {
			continue
		}

		proxiableContainers = append(proxiableContainers, ProxiableContainer{
			Name:        bridgeContainer.Name,
			Ipv4:        bridgeContainer.IPv4Address,
			VirtualHost: virtualHost,
			VirtualPort: virtualPort,
			Container:   &inspectedContainer,
		})
	}

	return proxiableContainers, nil
}

func RemoveContainer(name string) (string, error) {
	containers, _ := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{})

	var currentContainer types.Container

	for _, container := range containers {
		if container.Names[0] == "/"+name {
			currentContainer = container
		}
	}

	if currentContainer.ID == "" {
		return "", nil
	}

	err := GetDockerClient().ContainerStop(context.Background(), currentContainer.ID, nil)

	if err != nil {
		return "", err
	}

	err = GetDockerClient().ContainerRemove(context.Background(), currentContainer.Names[0], types.ContainerRemoveOptions{})

	if err != nil {
		return "", err
	}

	return containers[0].ID, nil
}
