package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/redwebcreation/hez/globals"
	"github.com/redwebcreation/hez/ui"
)

type RegistryAuth struct {
	Username string
	Password string
}

func (r RegistryAuth) ToBase64() string {
	auth, err := json.Marshal(map[string]string{
		"username": r.Username,
		"password": r.Password,
	})
	ui.Check(err)
	return base64.StdEncoding.EncodeToString(auth)
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
	Host         string
	Port         string
}

func (deployment DeploymentConfiguration) Deploy() {
	ui.Title("    " + deployment.Host)
	ui.NewLog("pulling %s", deployment.Image).Print()
	deployment.pullImage()
	ui.NewLog("successfully downloaded %s", deployment.Image).Print()

	ui.NewLog("stopping container %s", deployment.Name).Print()
	deployment.stopContainer()
	ui.NewLog("stopped container %s", deployment.Name).Top(1).Print()
	ui.NewLog("creating container %s", deployment.Name).Print()
	deployment.createContainer()
	ui.NewLog("created container %s", deployment.Name).Top(1).Print()

	if deployment.Healthchecks {
		ui.NewLog("checking the container healthyness").Print()
		deployment.waitForContainerToBeHealthy()
	} else {
		ui.NewLog("skipping healthchecks").Arrow(ui.Gray).Color(ui.Gray).ArrowString(" - ").Print()
	}

	if deployment.Warm {
		ui.NewLog("warming up server").Print()
		deployment.warmServer()
		ui.NewLog("server warmed up (went down from 100ms to 10ms)").Top(1).Print()
	} else {
		ui.NewLog("skipping server warmup").Arrow(ui.Gray).Color(ui.Gray).ArrowString(" - ").Print()
	}
}

func (d DeploymentConfiguration) pullImage() {
	pullOptions := types.ImagePullOptions{}

	if d.Registry != nil {
		pullOptions.RegistryAuth = d.Registry.ToBase64()
	}

	events, err := globals.Docker.ImagePull(context.Background(), d.Image, pullOptions)
	ui.Check(err)

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

	progress := ui.Progress{
		Prefix: "    " + "    " + ui.Bold + ui.White.AsFg(),
	}

	fmt.Println()
	progress.Render()

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			ui.Check(err)
		}

		progress.
			Increment(1).
			WithSuffix(ui.Bold + ui.Gray.AsFg() + strings.Replace(strings.ToLower(event.Status), "status: ", "", 1) + ui.Stop)
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
	ui.Check(err)

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
	ui.Check(err)

	err = globals.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
	ui.Check(err)

	if network != nil {
		// Force (re)connect
		_ = globals.Docker.NetworkDisconnect(context.Background(), network.ID, ref.ID, true)
		err = globals.Docker.NetworkConnect(context.Background(), network.ID, ref.ID, nil)
		ui.Check(err)
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
	ui.Check(err)

	var net types.NetworkResource

	if len(networks) > 0 {
		net, err = globals.Docker.NetworkInspect(context.Background(), networks[0].ID, types.NetworkInspectOptions{})
		ui.Check(err)
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
	variables := make([]string, len(d.Environment))

	for k, v := range d.Environment {
		variables = append(variables, k+"="+v)
	}

	return variables
}

func (d DeploymentConfiguration) waitForContainerToBeHealthy() {
	c := d.getContainer()

	inspection, err := globals.Docker.ContainerInspect(context.Background(), c.ID)
	ui.Check(err)

	isStarting := func(c types.ContainerJSON) bool {
		if c.State.Health == nil {
			return false
		}
		return c.State.Health.Status == "starting"
	}

	progress := ui.Progress{}
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
