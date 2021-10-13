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
	internal.Title("[" + deployment.Name + "]")
	internal.NewLog("pulling %s", deployment.Image).Print()
	deployment.pullImage()
	internal.NewLog("successfully downloaded %s", deployment.Image).Print()

	internal.NewLog("stopping container %s", deployment.Name).Print()
	deployment.stopContainer()
	internal.NewLog("stopped container %s", deployment.Name).Top(1).Print()
	internal.NewLog("creating container %s", deployment.Name).Print()
	deployment.createContainer()
	internal.NewLog("created container %s", deployment.Name).Top(1).Print()

	if deployment.Healthchecks {
		internal.NewLog("checking the container healthyness").Print()
		deployment.waitForContainerToBeHealthy()
	} else {
		internal.NewLog("skipping healthchecks").Arrow(internal.Gray).Color(internal.Gray).ArrowString(" - ").Print()
	}

	if deployment.Warm {
		internal.NewLog("warming up server").Print()
		deployment.warmServer()
		internal.NewLog("server warmed up (went down from 100ms to 10ms)").Top(1).Print()
	} else {
		internal.NewLog("skipping server warmup").Arrow(internal.Gray).Color(internal.Gray).ArrowString(" - ").Print()
	}
}

func (d DeploymentConfiguration) pullImage() {
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

	progress := internal.Progress{
		Prefix: "    " + "    " + internal.Bold + internal.Gray.AsFg(),
	}

	progress.Render()

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			internal.Check(err)
		}

		progress.
			Increment(1).
			WithSuffix(internal.Bold + internal.Gray.AsFg() + strings.ToLower(event.Status) + internal.Stop)
	}

	progress.Finish()
	fmt.Println()
}

func (d DeploymentConfiguration) stopContainer() *types.Container {
	c := d.getContainer()

	if c != nil {
		_ = globals.Docker.ContainerStop(context.Background(), c.ID, nil)
	}

	_ = globals.Docker.ContainerRemove(context.Background(), d.Name, types.ContainerRemoveOptions{})

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
			Name: "always",
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
	if d.Network == "" {
		return nil
	}

	networks, err := globals.Docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
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

	// TODO:
	if inspection.State.Health == nil || inspection.State.Health.Status == "healthy" {
		fmt.Println("healthy")
		return
	}

	fmt.Println("UNHEALTHY")
}

func (d DeploymentConfiguration) warmServer() {
	// TODO: Server warmup
}
