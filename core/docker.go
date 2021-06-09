package core

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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

func FindNetwork(name string) (types.NetworkResource, error) {
	networks, err := GetDockerClient().NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return types.NetworkResource{}, err
	}

	var network types.NetworkResource

	for _, maybeNetwork := range networks {
		if maybeNetwork.Name == name {
			network = maybeNetwork
			break
		}
	}

	networkDetails, err := GetDockerClient().NetworkInspect(context.Background(), network.ID, types.NetworkInspectOptions{})

	if err != nil {
		return networkDetails, err
	}

	return networkDetails, nil
}

func GetProxiableContainers() ([]ProxiableContainer, error) {
	config, _ := FindConfig(ConfigFile()).Resolve()

	networkDetails, err := FindNetwork(config.Network)

	if err != nil {
		return nil, err
	}

	var proxiableContainers []ProxiableContainer

	for _, networkContainer := range networkDetails.Containers {
		containersList, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: networkContainer.Name,
		})})

		if err != nil {
			return nil, err
		}

		container := containersList[0]

		inspectedContainer, err := GetDockerClient().ContainerInspect(context.Background(), container.ID)

		if err != nil {
			return nil, err
		}

		virtualHost := ""
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
			Name:        networkContainer.Name,
			Ipv4:        networkContainer.IPv4Address,
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
