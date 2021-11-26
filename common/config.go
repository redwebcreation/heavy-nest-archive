package common

import (
	"encoding/json"
	"github.com/docker/docker/pkg/homedir"
	"github.com/wormable/nest/ansi"
	"log/syslog"
	"os"
	"strings"
)

var ConfigFile string
var DataDirectory string
var CertificateDirectory string

func init() {
	configDirectory, err := homedir.GetConfigHome()
	ansi.Check(err)

	ConfigFile = configDirectory + "/nest/config.json"
	DataDirectory = configDirectory + "/nest/data"
	CertificateDirectory = configDirectory + "/nest/certs"

	err = os.MkdirAll(configDirectory+"/nest/data", os.FileMode(0700))
	ansi.Check(err)

	err = os.MkdirAll(configDirectory+"/nest/certs", os.FileMode(0700))
	ansi.Check(err)
}

type Configuration struct {
	DefaultNetwork     string `json:"default_network,omitempty"`
	DefaultMemoryLimit string `json:"default_memory_limit,omitempty"`

	Applications map[string]Application `json:"applications,omitempty"`

	Proxy struct {
		Logging          string `json:"logging,omitempty"`
		HttpPort         string `json:"http_port,omitempty"`
		HttpsPort        string `json:"https_port,omitempty"`
		SelfSigned       bool   `json:"self_signed,omitempty"`
	} `json:"proxy,omitempty"`

	ApiHost string `json:"api_host,omitempty"`

	Registries  []RegistryConfiguration `json:"registries,omitempty"`
	LogPolicies []LogPolicy             `json:"log_policies,omitempty"`
}

var Config Configuration

func LoadConfig() {
	_, err := os.Stat(ConfigFile)
	ansi.Check(err)

	contents, err := os.ReadFile(ConfigFile)
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

		if a.Quotas.Memory == "" {
			a.Quotas.Memory = config.DefaultMemoryLimit
		}

		if a.Quotas.CPU != 0 {
			a.Quotas.CPU = a.Quotas.CPU * 1000000000
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
