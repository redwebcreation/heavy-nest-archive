package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/redwebcreation/nest/cmd/ui"
)

type BackendStrategy string

type Configuration struct {
	DefaultNetwork string `json:"default_network,omitempty"`

	// an ip to a configuration
	// public  ips (disallowed, non-recommended?)
	Backends []string `json:"backends,omitempty"`

	Applications map[string]Application `json:"applications,omitempty"`

	Staging struct {
		Enabled bool   `json:"enabled,omitempty"`
		Host    string `json:"host,omitempty"`
		Logging string `json:"logging,omitempty"`

		MaxVersions int `json:"max_versions,omitempty"` // -1 for every commit, n for last n commits available in stating

		Database struct {
			Internal bool `json:"internal,omitempty"`

			Type string `json:"type,omitempty"`
			DSN  string `json:"dsn,omitempty"`
		} `json:"database"`

		Applications []string `json:"applications,omitempty"`
	} `json:"staging"`

	Production struct {
		Logging   string `json:"logging,omitempty"`
		HttpPort  string `json:"http_port,omitempty"`
		HttpsPort string `json:"https_port,omitempty"`
	} `json:"production"`

	BackendsManager struct {
		Host string `json:"host,omitempty"`
	}

	Registries  map[string]RegistryConfiguration `json:"registries,omitempty"`
	LogPolicies map[string][]LogPolicy           `json:"log_policies,omitempty"`
}

var Config Configuration

func init() {
	configFile := "/etc/nest/config.json"
	_, err := os.Stat(configFile)

	if err != nil {
		if os.IsNotExist(err) {
			ui.Check(
				fmt.Errorf("no config file found at %s", configFile),
			)
			return
		}

		ui.Check(err)
	}

	contents, err := os.ReadFile(configFile)
	ui.Check(err)

	Config = parseJsonConfig(contents)
}

func parseJsonConfig(contents []byte) Configuration {
	config := Configuration{
		DefaultNetwork: "bridge",
	}

	err := json.Unmarshal(contents, &config)
	ui.Check(err)

	applications := make(map[string]Application, len(config.Applications))

	for host, a := range config.Applications {
		a.Host = host

		if a.Port == "" {
			a.Port = "80"
		}

		if a.Network == "" {
			a.Network = config.DefaultNetwork
		}

		applications[host] = a
	}

	config.Applications = applications

	return config
}
