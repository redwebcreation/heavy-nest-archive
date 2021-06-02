package core

import (
	"crypto/sha256"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
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
		Redirections []struct {
			For   string `yaml:"for"`
			Kind  string `yaml:"kind"`
			Value string `yaml:"value"`
		}
	}
}

func EnsureConfigIsValid() {
	shouldExist := []string{
		ConfigDirectory(),
		ConfigFile(),
		ConfigDirectory() + "/compiled",
		ConfigDirectory() + "/environments",
		ConfigDirectory() + "/ssl",
	}

	var errors []error

	for _, requiredFile := range shouldExist {
		err := ensureFileExists(requiredFile)

		if err != nil {
			errors = append(errors, err)
		}

	}

	config, err := GetConfig()

	if err != nil {
		errors = append(errors, err)
	} else {
		for _, application := range config.Applications {
			_, err := os.Stat(ConfigDirectory() + "/environments/" + application.Environment)

			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		for _, configError := range errors {
			fmt.Println(configError)
		}

		fmt.Println("Configuration invalid.")
		os.Exit(1)
	}
}

func ensureFileExists(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return err
	}

	return nil
}

func GetConfig() (config Config, err error) {
	config = Config{}

	bytes, err := os.ReadFile(ConfigFile())
	data := string(bytes)

	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(data), &config)

	if err != nil {
		return config, err
	}

	return config, nil
}

func GetConfigChecksum() string {
	files := []string{ConfigFile()}
	checksum := ""

	filepath.Walk(ConfigDirectory(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "compiled") || strings.Contains(path, "ssl") || info.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	})

	for _, file := range files {
		checksum += getChecksumForFile(file)
	}

	return getChecksumForString(checksum)
}

func getChecksumForFile(file string) string {
	contents, err := os.ReadFile(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return getChecksumForString(string(contents))
}

func getChecksumForString(contents string) string {
	hash := sha256.New()
	input := strings.NewReader(contents)

	if _, err := io.Copy(hash, input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(hash.Sum(nil))
}
