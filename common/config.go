package common

import (
	"encoding/json"
	"log/syslog"
	"os"

	"github.com/wormable/ui"
)

type Configuration struct {
	DefaultNetwork string `json:"default_network,omitempty"`

	Applications map[string]Application `json:"applications,omitempty"`

	Production struct {
		Logging    string `json:"logging,omitempty"`
		HttpPort   string `json:"http_port,omitempty"`
		HttpsPort  string `json:"https_port,omitempty"`
		CertificateCache string `json:"certificate_cache,omitempty"`
		SelfSigned bool   `json:"self_signed,omitempty"`
	} `json:"production"`

	Registries  []RegistryConfiguration `json:"registries,omitempty"`
	LogPolicies []LogPolicy             `json:"log_policies,omitempty"`
}

var Config Configuration

func init() {
	configFile := "/etc/nest/config.json"
	_, err := os.Stat(configFile)

	if err != nil {
		if os.IsNotExist(err) {
			// TODO: Add it back once nestd is in a separate project
			//ui.Check(
			//	fmt.Errorf("no config file found at %s", configFile),
			//)
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

func (c Configuration) Log(level syslog.Priority, message string) {
	for _, policy := range c.LogPolicies {
		policy.Log(level, message)
	}
}
