package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/wormable/nest/globals"
	"strings"
)

type Service struct {
	Host     string
	Image    string   `json:"image" yaml:"image"`
	Mounts   []string `json:"mounts" yaml:"mounts"`
	Env      []string `json:"env" yaml:"env"`
	Prestart []string `json:"prestart" yaml:"prestart"`
}

func (s Service) Deploy() error {
	nextContainerName := s.NextContainerName()
	runningContainerName := s.RunningContainerName()

	err := StopContainerByName(nextContainerName)
	if err != nil {
		return err
	}

	c, err := s.createContainer(nextContainerName)
	if err != nil {
		return err
	}

	healthy, err := s.isHealthy(c)
	if err != nil {
		return err
	}

	if healthy {
		err = StopContainerByName(runningContainerName)
		if err != nil {
			return err
		}

		err = RenameContainer(nextContainerName, runningContainerName)
		if err != nil {
			return err
		}
	} else {
		// Trying to clean up, it's okay if it fails.
		_ = StopContainerByName(nextContainerName)

		return fmt.Errorf("container %s is not healthy", nextContainerName)
	}

	return nil
}

func (s Service) NextContainerName() string {
	return "next_" + strings.ReplaceAll(s.Host, ".", "_")
}

func (s Service) RunningContainerName() string {
	return strings.ReplaceAll(s.Host, ".", "_")
}

func (s Service) createContainer(name string) (string, error) {
	ref, err := globals.Docker.ContainerCreate(context.Background(), &container.Config{
		Image: s.Image,
	}, nil, nil, nil, name)
	if err != nil {
		return "", err
	}

	err = globals.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return ref.ID, nil
}

func (s Service) isHealthy(id string) (bool, error) {
	c, err := GetContainerById(id)
	if err != nil {
		return false, err
	}

	if c.Config.Healthcheck == nil {
		return true, nil
	}

	return c.State.Status == "healthy", nil
}
