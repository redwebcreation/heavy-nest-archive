package client

import (
	"strings"
)

type Application struct {
	Image     string            `json:"image,omitempty"`
	Host      string            `json:"host,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	EnvFiles  []string          `json:"env_files,omitempty"`
	Volumes   []string          `json:"volumes,omitempty"`
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

func (a Application) PrimaryContainer() DeploymentConfiguration {
	return DeploymentConfiguration{
		Image:       a.Image,
		Registry:    a.GetRegistry(),
		Environment: a.Env,
		Volumes:     a.Volumes,
		Network:     a.Network,
		Name:        strings.ReplaceAll(a.Host, ".", "_") + "_" + a.Port,
		Warm:        a.Warm,
		Host:        a.Host,
		Port:        a.Port,
	}
}

func (a Application) SecondaryContainer() DeploymentConfiguration {
	primary := a.PrimaryContainer()
	primary.Name += "_temporary"
	return primary
}
