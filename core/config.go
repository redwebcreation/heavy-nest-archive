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

type Config struct {
	Proxy        Proxy
	Applications []Application
	Logs         struct {
		Level        int8 `yaml:"level"`
		Redirections []string
	}
}

func IsConfigValid() error {
	_, err := os.Stat(ConfigFile())

	if os.IsNotExist(err) {
		return err
	}

	return nil
}

func EnsureConfigIsValid() {
	err := IsConfigValid()

	if err != nil {
		fmt.Println(err)
		fmt.Println("Configuration invalid.")
		os.Exit(1)
	}
}

func GetConfig() Config {
	config := Config{}

	bytes, _ := os.ReadFile(ConfigFile())

	err := yaml.Unmarshal(bytes, &config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return config
}

func GetConfigChecksum() string {
	return GetChecksumForFile(ConfigFile())
}

func GetChecksumForFile(file string) string {
	contents, err := os.ReadFile(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return GetChecksumForString(string(contents))
}

func GetChecksumForString(contents string) string {
	hash := sha256.New()
	input := strings.NewReader(contents)

	if _, err := io.Copy(hash, input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(hash.Sum(nil))
}
