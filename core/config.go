package core

import (
	"crypto/sha256"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

type Proxy struct {
	Port       string `yaml:"port"`
	Ssl        string `yaml:"ssl"`
	SelfSigned bool   `yaml:"self_signed"`
}

type ConfigData struct {
	Proxy        Proxy
	Applications []Application
	Logs         struct {
		Level        int8 `yaml:"level"`
		Redirections []string
	}
}

type Config string

func FindConfig(file string) Config {
	return Config(file)
}

func (config Config) IsValid() error {
	_, err := os.Stat(string(config))

	if os.IsNotExist(err) {
		return err
	}

	return nil
}

func (config Config) ExitIfInvalid() {
	err := config.IsValid()

	if err != nil {
		fmt.Println(err)
		fmt.Println("Configuration invalid.")
		os.Exit(1)
	}
}

func (config Config) Resolve() (ConfigData, error) {
	data := ConfigData{}
	bytes, _ := os.ReadFile(ConfigFile())

	err := yaml.Unmarshal(bytes, &config)

	if err != nil {
		return data, err
	}

	return data, nil
}

func (config Config) Checksum() (string, error) {
	return getChecksumForFile(string(config))
}

func getChecksumForFile(file string) (string, error) {
	contents, err := os.ReadFile(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hash := sha256.New()
	input := strings.NewReader(string(contents))

	if _, err := io.Copy(hash, input); err != nil {
		return "", err
	}

	return string(hash.Sum(nil)), nil
}
