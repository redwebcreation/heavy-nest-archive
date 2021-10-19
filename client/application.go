package client

type Application struct {
	Image     string            `json:"image,omitempty"`
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
