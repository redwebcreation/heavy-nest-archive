package core

import (
	"io/ioutil"
	"strings"
)

type Environment struct {
	Name     string `yaml:"name"`
	Filename string `yaml:"filename"`
}

func (environment Environment) Contents() string {
	bytes, err := ioutil.ReadFile(
		ConfigDirectory() + "/envs/" + environment.Filename,
	)
	data := string(bytes)

	if err != nil {
		return ""
	}

	return strings.TrimSpace(data)
}

func FindEnvironment(name string) (Environment, error) {
	config, err := GetConfig()

	if err != nil {
		return Environment{}, err
	}

	for _, environment := range config.Environments {
		if environment.Name == name {
			return environment, nil
		}
	}

	return Environment{}, nil
}
