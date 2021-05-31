package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"os"
	"os/exec"
	"strings"
)

type Application struct {
	Name        string `yaml:"name"`
	Domain      string `yaml:"domain"`
	Image       string `yaml:"image"`
	Environment string `yaml:"environment"`
}

func (application Application) Start() (*Application, error) {
	dockerCommand := []string{"run", "-d", "--name", application.Name}

	environment, err := FindEnvironment(application.Environment)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(environment.Contents()), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		dockerCommand = append(dockerCommand, "-e", line)
	}

	dockerCommand = append(dockerCommand, "-e", "VIRTUAL_HOST="+application.Domain)
	dockerCommand = append(dockerCommand, application.Image)

	var stdErr bytes.Buffer
	cmd := exec.Command("docker", dockerCommand...)
	cmd.Stderr = &stdErr
	err = cmd.Run()

	if stdErr.Len() > 0 {
		fmt.Println(stdErr.String())
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &application, nil
}

func (application Application) GetEnvironment() (Environment, error) {
	return FindEnvironment(application.Environment)
}

func (application Application) HasRunningContainer() bool {
	return application.GetContainer().ID != ""
}

func (application Application) GetContainer() types.Container {
	containers, err := GetDockerClient().ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, container := range containers {
		if container.Image != application.Image {
			continue
		}

		for _, name := range container.Names {
			if "/"+application.Name == name {
				return container
			}
		}
	}

	return types.Container{}
}

func (application Application) StopContainer() *Application {
	if !application.HasRunningContainer() {
		return &application
	}

	err := GetDockerClient().ContainerStop(context.Background(), application.GetContainer().ID, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &application
}
