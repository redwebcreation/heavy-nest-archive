package internal

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

	data.useDefaults()

	Config = &data
}

func (h *Hez) useDefaults() {
	if h.DefaultNetwork == "" {
		h.DefaultNetwork = "bridge"
	}

	if h.Proxy.Http.Port == "" {
		h.Proxy.Http.Port = "80"
	}

	if h.Proxy.Https.Port == "" {
		h.Proxy.Http.Port = "443"
	}

	if h.Proxy.Https.SelfSigned == nil {
		selfSigned := false
		h.Proxy.Https.SelfSigned = &selfSigned
	}

	for i := range h.Applications {
		h.Applications[i].Host = i

		if h.Applications[i].ContainerPort == "" {
			h.Applications[i].ContainerPort = "80"
		}

		if h.Applications[i].Network == "" {
			h.Applications[i].Network = h.DefaultNetwork
		}

		if h.Applications[i].Warm == nil {
			warm := true
			h.Applications[i].Warm = &warm
		}
	}
}
