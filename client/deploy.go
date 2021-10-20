package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/redwebcreation/nest/globals"
	"github.com/redwebcreation/nest/ui"
)

type RegistryConfiguration struct {
	Host     string
	Username string
	Password string
}

func (r RegistryConfiguration) ToBase64() string {
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
	Image         string
	Registry      *RegistryConfiguration
	Environment   map[string]string
	Volumes       []Volume
	Network       string
	Name          string
	Warm          bool
	Host          string
	Port          string
}
type DeploymentOptions struct {
	Pull         bool
	Healthchecks bool
}

func (d DeploymentConfiguration) Deploy(opts DeploymentOptions) {
	if opts.Pull {
		ui.NewLog("pulling %s", d.Image).Print()
		d.pullImage()
		ui.NewLog("successfully downloaded %s", d.Image).Print()
	}

	ui.NewLog("stopping container %s", d.Name).Print()
	stopped := d.StopContainer()

	if stopped {
		ui.NewLog("stopped container %s", d.Name).Top(1).Print()
	} else {
		fmt.Println("\033[3A\033[K") // We erase the "stopping container" log
	}

	ui.NewLog("creating container %s", d.Name).Print()
	d.createContainer()
	ui.NewLog("created container %s", d.Name).Top(1).Print()

	if opts.Healthchecks {
		ui.NewLog("checking the container healthyness").Print()
		d.waitForContainerToBeHealthy()
	} else {
		ui.NewLog("skipping healthchecks").Arrow(ui.Gray).Color(ui.Gray).ArrowString(" - ").Print()
	}

	if d.Warm {
		ui.NewLog("warming up server").Print()
		d.warmServer()
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
		Total:  60,
		Prefix: "    " + "    " + ui.White.Fg(),
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
			WithSuffix(ui.Gray.Fg() + strings.Replace(strings.ToLower(event.Status), "status: ", "", 1) + ui.Stop)
	}

	progress.Finish()
	fmt.Println()
}

func (d DeploymentConfiguration) StopContainer() bool {
	c := d.getContainer()

	if c == nil {
		return false
	}

	err := globals.Docker.ContainerStop(context.Background(), c.ID, nil)
	ui.Check(err)

	err = globals.Docker.ContainerRemove(context.Background(), d.Name, types.ContainerRemoveOptions{})
	ui.Check(err)

	return true
}

func (d DeploymentConfiguration) getContainer() *types.Container {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "name", Value: d.Name}),
		All: true,
	})
	ui.Check(err)

	for _, c := range containers {
		// Filters return non-exact matches so {Key: name, Value: example_com_80}
		// could return example_com_80,example_com_80_temporary...qqq
		if c.Names[0] == "/"+d.Name {
			return &c
		}
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

	if d.Network != "" {
		// Force (re)connect
		_ = globals.Docker.NetworkDisconnect(context.Background(), network.ID, ref.ID, true)
		err = globals.Docker.NetworkConnect(context.Background(), network.ID, ref.ID, nil)
		ui.Check(err)
	}

	return ref.ID
}

func (d DeploymentConfiguration) getNetwork() types.NetworkResource {
	networks, err := globals.Docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: d.Network,
		}),
	})
	ui.Check(err)

	// We don't have to check if the network exists, it does, we check for it before we even run the command.
	net, err := globals.Docker.NetworkInspect(context.Background(), networks[0].ID, types.NetworkInspectOptions{})
	ui.Check(err)
	return net
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
	variables := make([]string, len(d.Environment)-1)

	for k, v := range d.Environment {
		variables = append(variables, k+"="+v)
	}

	return variables
}

func (d DeploymentConfiguration) waitForContainerToBeHealthy() {
	c := d.getContainer()

	inspection, err := globals.Docker.ContainerInspect(context.Background(), c.ID)
	ui.Check(err)

	if inspection.State.Health == nil {
		ui.NewLog("no healthchecks defined").Arrow(ui.Gray).ArrowString("    " + "- ").Color(ui.Gray).Print()
		return
	}

	progress := ui.Progress{
		Total: int(inspection.Config.Healthcheck.Interval.Seconds()),
	}

	progress.Render()
	for inspection.State.Health.Status == "starting" {
		progress.Increment(1)
		time.Sleep(time.Second)
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
	if d.Network == "" {
		d.getNetwork()
	}
}
