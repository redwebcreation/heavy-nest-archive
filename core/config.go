package core

import (
	"crypto/sha256"
	"encoding/hex"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

type Proxy struct {
	Port       int   `yaml:"port",omitempty`
	Ssl        int   `yaml:"ssl",omitempty`
	SelfSigned *bool `yaml:"self_signed",omitempty`
}

type ConfigData struct {
	Proxy        Proxy
	Applications []Application
	Logs         struct {
		Level        int8 `yaml:"level",omitempty`
		Redirections []string
	}
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

	return data
}

func (config Config) Checksum() (string, error) {
	contents, err := os.ReadFile(string(config))

	if err != nil {
		return "", err
	}

	hash := sha256.New()
	input := strings.NewReader(string(contents))

	if _, err := io.Copy(hash, input); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
