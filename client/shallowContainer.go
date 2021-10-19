package client

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/redwebcreation/nest/globals"
	"github.com/redwebcreation/nest/ui"
)

type ShallowContainer struct {
	Name        string
	Application Application
}

func (s ShallowContainer) get() *types.Container {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Limit: 1,
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "name",
				Value: s.Name,
			},
		),
	})
	ui.Check(err)

	if len(containers) == 0 {
		return nil
	}

	return &containers[0]
}

func (s ShallowContainer) Stop() {
	c := s.get()

	if c == nil || c.State != "running" {
		return
	}

	// an error is thrown error if the container dow
	err := globals.Docker.ContainerStop(context.Background(), c.ID, nil)
	ui.Check(err)
}

func (s ShallowContainer) Start() {
	registry := s.Application.GetRegistry()

	conf := DeploymentConfiguration{
		Image:       s.Application.Image,
		Environment: s.Application.Env,
		Volumes:     s.Application.Volumes,
		Network:     s.Application.Network,
		Name:        s.Name,
		Host:        s.Application.Host,
		Warm:        s.Application.Warm,
		Port:        s.Application.Port,
	}

	if registry != nil {
		conf.Registry = registry
	}

	conf.Deploy()
}
