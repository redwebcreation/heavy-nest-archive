package common

import (
	"encoding/json"
	"fmt"
	"github.com/wormable/nest/ansi"
	"log/syslog"
	"os"
	"strings"
)

type Configuration struct {
	DefaultNetwork     string `json:"default_network,omitempty"`
	DefaultMemoryLimit string `json:"default_memory,omitempty"`

	Applications map[string]Application `json:"applications,omitempty"`

	Proxy struct {
		Logging          string `json:"logging,omitempty"`
		HttpPort         string `json:"http_port,omitempty"`
		HttpsPort        string `json:"https_port,omitempty"`
		CertificateCache string `json:"certificate_cache,omitempty"`
		SelfSigned       bool   `json:"self_signed,omitempty"`
	} `json:"proxy,omitempty"`
	Registries  []RegistryConfiguration `json:"registries,omitempty"`
	LogPolicies []LogPolicy             `json:"log_policies,omitempty"`
}

var Config Configuration

func init() {
	LoadConfig()
}

func LoadConfig() {
	configFile := "/etc/nest/config.json"
	_, err := os.Stat(configFile)

	if err != nil {
		if os.IsNotExist(err) {
			ansi.Check(
				fmt.Errorf("no config file found at %s", configFile),
			)
			return
		}

		ansi.Check(err)
	}

	contents, err := os.ReadFile(configFile)
	ansi.Check(err)

	Config = parseJsonConfig(contents)
}

func parseJsonConfig(contents []byte) Configuration {
	config := Configuration{
		DefaultNetwork:     "bridge",
		DefaultMemoryLimit: "-1",
	}

	err := json.Unmarshal(contents, &config)
	ansi.Check(err)

	applications := make(map[string]Application)

	for host, a := range config.Applications {
		a.Host = host

		if a.Port == "" {
			a.Port = "80"
		}

		if a.Network == "" {
			a.Network = config.DefaultNetwork
		}

		if a.Registry != "" {
			// we can't use the GetRegistry method as it uses the global Config variable
			// that hasn't been set yet
			for _, registry := range config.Registries {
				if registry.Name == a.Registry {
					a.Image = strings.TrimRight(registry.Host, "/") + "/" + strings.TrimLeft(a.Image, "/")
				}
			}
		}

		if a.Quota.Memory == "" {
			a.Quota.Memory = config.DefaultMemoryLimit
		}

		applications[host] = a
	}

	config.Applications = applications

	for _, application := range config.Applications {
		for _, alias := range application.Aliases {
			config.Applications[alias] = application
		}
	}

	return config
}

func (c Configuration) Log(level syslog.Priority, message string) {
	for _, policy := range c.LogPolicies {
		policy.Log(level, message)
	}
}
