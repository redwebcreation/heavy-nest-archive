package core

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"os"
	"strings"
)

func GetDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

	for _, proxiableContainer := range GetProxiableContainers() {
		domains = append(domains, proxiableContainer.VirtualHost)
	}

	return domains
}

func GetProxiableContainers() []ProxiableContainer {
	networks, err := GetDockerClient().NetworkList(context.Background(), types.NetworkListOptions{})
	var proxiableContainers []ProxiableContainer

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var networkBridge types.NetworkResource

	for _, network := range networks {
		if network.Name == "bridge" {
			networkBridge = network
		}
	}

	bridgeDetails, err := GetDockerClient().NetworkInspect(context.Background(), networkBridge.ID, types.NetworkInspectOptions{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, bridgeContainer := range bridgeDetails.Containers {
		containerList, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(containerList) != 1 {
			fmt.Println("This should not happen. Zero or more than one container found.")
			os.Exit(1)
		}

		container := containerList[0]
		inspectedContainer, err := GetDockerClient().ContainerInspect(context.Background(), container.ID)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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

	return proxiableContainers
}
