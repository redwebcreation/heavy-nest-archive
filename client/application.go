package client

import (
	"fmt"
)

type Application struct {
	Image     string            `json:"image"`
	Env       map[string]string `json:"env"`
	EnvFiles  []string          `json:"env_files"`
	Volumes   []Volume          `json:"volumes"`
	Warm      bool              `json:"warm"`
	Backend   BackendStrategy   `json:"backend"`
	LogPolicy string            `json:"log_policy"`
	Registry  string            `json:"registry"`
	Network   string            `json:"network"`
	Port      string            `json:"port"`
}

func (a Application) GetRegistry() (*RegistryAuth, error) {
	if a.Registry == "" {
		return nil, nil
	}

	for name, registry := range Config.Registries {
		if name == a.Registry {
			return &registry, nil
		}
	}

	return nil, fmt.Errorf("registry [%s] not found", a.Registry)
}

func (a Application) GetDeploymentConfigurationFor(name string) (DeploymentConfiguration, error) {
	registry, err := a.GetRegistry()
	if err != nil {
		return DeploymentConfiguration{}, err
	}

	conf := DeploymentConfiguration{
		Image:       a.Image,
		Environment: a.Env,
		Volumes:     a.Volumes,
		Network:     a.Network,
		Name:        name,
		Warm:        a.Warm,
		Port:        a.Port,
	}

	if registry != nil {
		conf.Registry = registry
	}

	return conf, nil
}
