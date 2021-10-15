package client

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/redwebcreation/hez/ui"
)

type BackendStrategy string

type Configuration struct {
	DefaultNetwork string

	// an ip to a configuration
	// self can be used to use the current server as a backend
	// public  ips (disallowed, non-recommended?)
	Backends []string

	Applications map[string]Application

	Staging struct {
		Enabled bool
		Host    string
		Logging string

		MaxVersions int // -1 for every commit, n for last n commits available in stating

		Database struct {
			Internal bool

			Type string
			DSN  string
		}

		Applications []string
	}

	Production struct {
		Logging   string
		HttpPort  string
		HttpsPort string
	}

	Registries  map[string]RegistryAuth
	LogPolicies map[string]LogPolicy
}

var Config Configuration

func init() {
	configHome := "/etc/hez"
	configFiles := []string{"config.yml", "config.yaml", "config.json"}
	var configFile string
	for _, maybeConfigFile := range configFiles {
		_, err := os.Stat(configHome + "/" + maybeConfigFile)

		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			ui.Check(err)
		} else {
			configFile = configHome + "/" + maybeConfigFile
		}
	}

	if configFile == "" {
		ui.Check(fmt.Errorf("no config file found in %s/config.{yml,yaml,json}", configHome))
	}

	contents, err := os.ReadFile(configFile)
	ui.Check(err)

	Config = parseJsonConfig(contents)
}

func parseJsonConfig(contents []byte) Configuration {
	var config Configuration

	err := json.Unmarshal(contents, &config)
	ui.Check(err)
	return config
}
