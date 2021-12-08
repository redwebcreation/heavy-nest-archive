package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/wormable/nest/globals"
)

func GetContainer(name string) (*types.ContainerJSON, error) {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "name", Value: name}),
	})

	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("no container found matching %s", name)
	}

	if len(containers) > 1 {
		return nil, fmt.Errorf("ambiguous container name: multiple running containers found matching %s", name)
	}

	return InspectContainer(containers[0].ID)
}

func GetContainerById(id string) (*types.ContainerJSON, error) {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "id", Value: id},
		),
	})

	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("no container found with id %s", id)
	}

	// We assume that there is only one container with the given id
	// This assumption may be false if we pass in something like "c8"
	// and there are multiple containers with an id starting with "c8"
	// However, we never call this function with non-complete ids.
	// Therefore, it is a safe assumption to make.
	return InspectContainer(containers[0].ID)
}

func InspectContainer(id string) (*types.ContainerJSON, error) {
	inspection, err := globals.Docker.ContainerInspect(context.Background(), id)
	return &inspection, err
}

func StopContainerByName(name string) error {
	c, err := GetContainer(name)
	if err != nil {
		return err
	}

	return globals.Docker.ContainerStop(context.Background(), c.ID, nil)
}

func RenameContainer(id string, newName string) error {
	return globals.Docker.ContainerRename(context.Background(), id, newName)
}
