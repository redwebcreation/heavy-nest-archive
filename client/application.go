package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/redwebcreation/nest/globals"
	"github.com/redwebcreation/nest/ui"
)

type Application struct {
	Image     string            `json:"image,omitempty"`
	Host      string            `json:"host,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	EnvFiles  []string          `json:"env_files,omitempty"`
	Volumes   []Volume          `json:"volumes,omitempty"`
	Warm      bool              `json:"warm,omitempty"`
	Backend   BackendStrategy   `json:"backend,omitempty"`
	LogPolicy string            `json:"log_policy,omitempty"`
	Registry  string            `json:"registry,omitempty"`
	Network   string            `json:"network,omitempty"`
	Port      string            `json:"port,omitempty"`
}

func (a Application) GetRegistry() *RegistryConfiguration {
	if a.Registry == "" {
		return nil
	}

	for name, registry := range Config.Registries {
		if name == a.Registry {
			return &registry
		}
	}

	return nil
}

func (a Application) GetDeploymentConfigurationFor(name string, host string) DeploymentConfiguration {
	registry := a.GetRegistry()

	conf := DeploymentConfiguration{
		Image:       a.Image,
		Environment: a.Env,
		Volumes:     a.Volumes,
		Network:     a.Network,
		Name:        name,
		Host:        host,
		Warm:        a.Warm,
		Port:        a.Port,
	}

	if registry != nil {
		conf.Registry = registry
	}

	return conf
}

func (a Application) getContainerName() string {
	return strings.ReplaceAll(a.Host, ".", "_") + "_" + a.Port
}

type LazyContainer struct {
	Name string
}

func (a Application) PrimaryContainer() LazyContainer {
	return LazyContainer{
		Name: a.getContainerName(),
	}
}

func (a Application) SecondaryContainer() LazyContainer {
	return LazyContainer{
		Name: a.getContainerName() + "_temporary",
	}
}

func (c LazyContainer) get() *types.Container {
	c, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Limit:   1,
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key: "name",
				Value: c.Name,
			}
		),
	})
	ui.Check(err)

	if len(c)==0 {
		ui.Check(fmt.Errorf("container %s does not exist", c.Name))
	}

	return c[0]
}
func (c LazyContainer) Stop()  {}
func (c LazyContainer) Start() {}
