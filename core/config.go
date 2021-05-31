package core

import (
	"crypto/sha256"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Applications []Application
	Environments []Environment
	Logs         struct {
		MaxSize string `yaml:"max_size"`
		Beacon  struct {
			Url   string `yaml:"url"`
			Every string `yaml:"every"`
		}
	}
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
	//var checksums []string
	files := []string{ConfigFile()}
	checksum := ""

	filepath.Walk(configDirectory(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "compiled") || info.IsDir() {
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
		log.Fatal(err)
	}

	return getChecksumForString(string(contents))
}

func getChecksumForString(contents string) string {
	hash := sha256.New()
	input := strings.NewReader(contents)

	if _, err := io.Copy(hash, input); err != nil {
		log.Fatal(err)
	}

	return string(hash.Sum(nil))
}
