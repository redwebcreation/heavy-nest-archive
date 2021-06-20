package core

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Hez struct {
	DefaultNetwork string                  `yaml:"default_network"`
	Applications   map[string]*Application `yaml:"applications"`
	Proxy          struct {
		Logs struct {
			Level        int      `yaml:"level"`
			Redirections []string `yaml:"redirections"`
		} `yaml:"logs"`
		Http struct {
			Port string `yaml:"port"`
		} `yaml:"http"`
		Https struct {
			Port       string `yaml:"port"`
			SelfSigned *bool  `yaml:"self_signed"`
		} `yaml:"https"`
		Firewall struct {
			Rate int `yaml:"rate"` // max requests every 5 minutes per ip.
		} `yaml:"firewall"`
	} `yaml:"proxy"`
}

var Config *Hez
var ConfigFile = "/etc/hez/hez.yml"

func init() {
	data := Hez{}
	bytes, _ := os.ReadFile(ConfigFile)

	err := yaml.Unmarshal(bytes, &data)

	if err != nil {
		panic(err)
	}

	useDefaults(&data)

	Config = &data
}

func useDefaults(config *Hez) {
	if config.DefaultNetwork == "" {
		config.DefaultNetwork = "bridge"
	}

	if config.Proxy.Http.Port == "" {
		config.Proxy.Http.Port = "80"
	}

	if config.Proxy.Https.Port == "" {
		config.Proxy.Http.Port = "443"
	}

	if config.Proxy.Https.SelfSigned == nil {
		selfSigned := false
		config.Proxy.Https.SelfSigned = &selfSigned
	}

	for i := range config.Applications {
		if config.Applications[i].ContainerPort == "" {
			config.Applications[i].ContainerPort = "80"
		}

		if config.Applications[i].Network == "" {
			config.Applications[i].Network = config.DefaultNetwork
		}

		if config.Applications[i].Warm == nil {
			warm := true
			config.Applications[i].Warm = &warm
		}
	}
}
