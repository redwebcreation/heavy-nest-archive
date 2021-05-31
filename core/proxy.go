package core

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"os"
)

type ProxiableContainer struct {
	Name string
	ID   string
	Ip   string
}

func GetProxiableContainers() []types.EndpointResource {
	networks, err := GetDockerClient().NetworkList(context.Background(), types.NetworkListOptions{})
	var containers []types.EndpointResource

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
		containers = append(containers, bridgeContainer)
	}

	return containers
}
