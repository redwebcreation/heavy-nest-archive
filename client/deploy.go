package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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

type DeploymentConfiguration struct {
	Image       string
	Registry    *RegistryConfiguration
	Environment map[string]string
	Volumes     []string
	Network     string
	Name        string
	Warm        bool
	Host        string
	Port        string
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
		ui.NewLog("server warmed up").Print()
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
	net := d.getNetwork()

	// TODO: Quotas
	// TODO: Better volumes
	ref, err := globals.Docker.ContainerCreate(context.Background(), &container.Config{
		Env:   d.EnvironmentToDockerEnv(),
		Image: d.Image,
		Labels: map[string]string{
			"nest-id":   d.Host + d.Port,
			"nest-host": d.Host,
			"nest-port": d.Port,
		},
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Binds: d.Volumes,
	}, nil, nil, d.Name)
	ui.Check(err)

	err = globals.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
	ui.Check(err)

	// Docker's default network
	_ = globals.Docker.NetworkDisconnect(context.Background(), net.ID, ref.ID, true)
	err = globals.Docker.NetworkConnect(context.Background(), net.ID, ref.ID, nil)
	ui.Check(err)

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
	_ = d.getContainer().NetworkSettings.Networks[d.Network]

	var responseTimes [10]float64

	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := http.Get("http://172.17.0.3")
		end := time.Now()
		ui.Check(err)

		responseTime := end.Add(time.Duration(-start.Nanosecond()))

		responseTimes[i] = float64(responseTime.Nanosecond()) / float64(time.Millisecond)
		responseTimeString := fmt.Sprintf("%.2fms", responseTimes[i])

		fmt.Printf(
			"      %sGET / in%s %s\n",
			ui.Gray.Fg(),
			ui.Stop,
			ui.Primary.Fg()+responseTimeString+ui.Stop,
		)
	}

	var max float64
	var min *float64
	var sum float64
	var correction float64
	diff := responseTimes[0] - responseTimes[9]

	for _, responseTime := range responseTimes {
		if responseTime > max {
			max = responseTime
		}

		if min == nil || responseTime < *min {
			min = &responseTime
		}

		sum += responseTime

		correction += (diff) - (responseTimes[0] - responseTime)
	}

	diff -= correction / 10.0

	fmt.Println()
	if diff <= 0.02 {
		fmt.Print(ui.Yellow.Fg())
		fmt.Println("      It seems like your server does not need to warmup.")
		fmt.Println("      Consider removing this option in your config to speed up the process.")
		fmt.Print(ui.Stop)
	}

	fmt.Printf("      %sGain: %.2fms%s Diff: %.2fms Max: %.2fms Min: %.2fms Avg: %.2fms\n\n", ui.White.Fg(), diff, ui.Stop, responseTimes[0]-responseTimes[9], max, *min, sum/10.0)
}
