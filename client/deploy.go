package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/redwebcreation/hez/globals"
	"github.com/redwebcreation/hez/internal"
)

type RegistryAuth struct {
	Username string
	Password string
}

type Volume struct {
	From string
	To   string
}

type DeploymentConfiguration struct {
	Image        string
	Registry     *RegistryAuth
	Environment  map[string]string
	Volumes      []Volume
	Network      string
	Name         string
	Warm         bool
	Healthchecks bool
}

func (deployment DeploymentConfiguration) Deploy() {
	if !deployment.hasLatestImage() {
		deployment.pullImage()
	}

	deployment.stopContainer()
	deployment.createContainer()

	if deployment.Healthchecks {
		deployment.waitForContainerToBeHealthy()
	}

	if deployment.Warm {
		deployment.warmServer()
	}
}

func (d DeploymentConfiguration) hasLatestImage() bool {
	images, err := globals.Docker.ImageList(context.Background(), types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: d.Image,
		},
		),
	})
	internal.Check(err)

	fmt.Println(images)
	return true
}

func (d DeploymentConfiguration) pullImage() error {
	pullOptions := types.ImagePullOptions{}

	if d.Registry != nil {
		auth, err := json.Marshal(map[string]string{
			"username": d.Registry.Username,
			"password": d.Registry.Password,
		})
		internal.Check(err)
		pullOptions.RegistryAuth = base64.StdEncoding.EncodeToString(auth)
	}

	events, err := globals.Docker.ImagePull(context.Background(), d.Image, pullOptions)
	internal.Check(err)

	decoder := json.NewDecoder(events)

	var event *struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	progress := internal.Progress{}
	progress.Render()

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		progress.Increment(1).WithLabel(event.Status)
	}

	progress.Finish()
	return nil
}

func (d DeploymentConfiguration) stopContainer() *types.Container {
	c := d.getContainer()

	if c != nil {
		_ = globals.Docker.ContainerStop(context.Background(), c.ID, nil)
	}

	err := globals.Docker.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{})
	internal.Check(err)

	return c
}

func (d DeploymentConfiguration) getContainer() *types.Container {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Limit: 1,
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "name",
				Value: d.Name,
			},
		),
	})
	internal.Check(err)

	if len(containers) > 0 {
		return &containers[0]
	}

	return nil
}

func (d DeploymentConfiguration) createContainer() string {
	network := d.getNetwork()

	ref, err := globals.Docker.ContainerCreate(context.Background(), &container.Config{
		Env:   d.EnvironmentToDockerEnv(),
		Image: d.Image,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name:              "always",
			MaximumRetryCount: 2,
		},
		Mounts: d.VolumesToDockerMounts(),
	}, nil, nil, d.Name)
	internal.Check(err)

	err = globals.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
	internal.Check(err)

	if network != nil {
		// Force (re)connect
		_ = globals.Docker.NetworkDisconnect(context.Background(), network.ID, ref.ID, true)
		err = globals.Docker.NetworkConnect(context.Background(), network.ID, ref.ID, nil)
		internal.Check(err)
	}

	return ref.ID
}

func (d DeploymentConfiguration) getNetwork() *types.NetworkResource {
	networks, err := globals.Docker.NetworkList(context.Background(), types.NetworkListOptions{
		filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: d.Network,
		}),
	})
	internal.Check(err)

	var net types.NetworkResource

	if len(networks) > 0 {
		net, err = globals.Docker.NetworkInspect(context.Background(), networks[0].ID, types.NetworkInspectOptions{})
		internal.Check(err)
		return &net
	}
	return nil
}

func (d DeploymentConfiguration) VolumesToDockerMounts() []mount.Mount {
	var mounts []mount.Mount

	for _, volume := range d.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: volume.From,
			Target: volume.To,
		})
	}

	return mounts
}

func (d DeploymentConfiguration) EnvironmentToDockerEnv() []string {
	var variables []string

	for k, v := range d.Environment {
		if v == "" && !strings.Contains(k, "=") {
			contents, err := os.ReadFile(k)
			internal.Check(err)

			for _, envFileVariable := range strings.Split(string(contents), "\n") {
				trimmed := strings.TrimSpace(envFileVariable)

				if trimmed == "" {
					continue
				}

				variables = append(variables, trimmed)
			}

			continue
		}

		variables = append(variables, k+"="+v)
	}

	return variables
}

func (d DeploymentConfiguration) waitForContainerToBeHealthy() {
	c := d.getContainer()

	inspection, err := globals.Docker.ContainerInspect(context.Background(), c.ID)
	internal.Check(err)

	isStarting := func(c types.ContainerJSON) bool {
		if c.State.Health == nil {
			return false
		}
		return c.State.Health.Status == "starting"
	}

	progress := internal.Progress{}
	progress.Render()
	for isStarting(inspection) {
		progress.Increment(1)
	}
	progress.Finish()

	if inspection.State.Health == nil || inspection.State.Health.Status == "healthy" {
		fmt.Println("healthy")
		return
	}
}

func (d DeploymentConfiguration) warmServer() {
	// TODO: Server warmup
}
