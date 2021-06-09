package core

import (
	"gopkg.in/yaml.v2"
	"os"
)

type ConfigData struct {
	Network string `yaml:"network,omitempty"`
	Proxy struct {
		Port       int   `yaml:"port,omitempty"`
		Ssl        int   `yaml:"ssl,omitempty"`
		SelfSigned *bool `yaml:"self_signed,omitempty"`
		Logs       struct {
			Level        int8 `yaml:"level,omitempty"`
			Redirections []string
		}
	}
	Applications []Application
}

type Config string

func FindConfig(file string) Config {
	return Config(file)
}

func (config Config) IsValid() bool {
	_, err := os.Stat(string(config))

	if os.IsNotExist(err) {
		return false
	}

	resolved, err := config.Resolve()

	if err != nil {
		return false
	}

	var usedHosts = make([]string, len(resolved.Applications))
	for _, application := range resolved.Applications {
		for _, host := range usedHosts {
			if host == application.Host {
				return false
			}
		}

		usedHosts = append(usedHosts, application.Host)
	}

	return true
}
func (config Config) Resolve() (ConfigData, error) {
	data := ConfigData{}
	bytes, _ := os.ReadFile(string(config))

	err := yaml.Unmarshal(bytes, &data)

	if err != nil {
		return data, err
	}

	useDefaults(&data)

	return data, nil
}

func useDefaults(data *ConfigData) *ConfigData {
	if data.Network == "" {
		data.Network = "bridge"
	}
	if data.Proxy.Port == 0 {
		data.Proxy.Port = 80
	}

	if data.Proxy.Ssl == 0 {
		data.Proxy.Ssl = 443
	}

	if data.Proxy.SelfSigned == nil {
		selfSigned := false
		data.Proxy.SelfSigned = &selfSigned
	}

	for i := range data.Applications {
		if data.Applications[i].ContainerPort == 0 {
			data.Applications[i].ContainerPort = 80
		}
	}

	return data
}