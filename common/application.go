package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/wormable/ui"
	"github.com/wormable/nest/globals"
)

type Application struct {
	Image     string            `json:"image,omitempty"`
	Host      string            `json:"host,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	EnvFiles  []string          `json:"env_files,omitempty"`
	Volumes   []string          `json:"volumes,omitempty"`
	Warm      bool              `json:"warm,omitempty"`
	Backend   BackendStrategy   `json:"backend,omitempty"`
	LogPolicy string            `json:"log_policy,omitempty"`
	Registry  string            `json:"registry,omitempty"`
	Network   string            `json:"network,omitempty"`
	Port      string            `json:"port,omitempty"`
}

func (a Application) GetRegistry() *RegistryConfiguration {
	if a.Registry == "" {
		return nil
	}

	for name, registry := range Config.Registries {
		if name == a.Registry {
			return &registry
		}
	}

	return nil
}

type DeploymentOptions struct {
	Pull         bool
	Healthchecks bool
	Name         string
}

func (a Application) Deploy(opts DeploymentOptions) {
	if opts.Pull {
		ui.NewLog("pulling %s", a.Image).Print()
		a.pullImage()
		ui.NewLog("successfully downloaded %s", a.Image).Print()
	}

	ui.NewLog("stopping container %s", opts.Name).Print()
	stopped := a.StopContainer(opts.Name)

	if stopped {
		ui.NewLog("stopped container %s", opts.Name).Top(1).Print()
	} else {
		fmt.Println("\033[3A\033[K") // We erase the "stopping container" log
	}

	ui.NewLog("creating container %s", opts.Name).Print()
	a.createContainer(opts.Name)
	ui.NewLog("created container %s", opts.Name).Top(1).Print()

	if opts.Healthchecks {
		ui.NewLog("checking the container healthyness").Print()
		a.waitForContainerToBeHealthy(opts.Name)
	} else {
		ui.NewLog("skipping healthchecks").Arrow(ui.Gray).Color(ui.Gray).ArrowString(" - ").Print()
	}

	if a.Warm {
		ui.NewLog("warming up server").Print()
		a.warmServer(opts.Name)
		ui.NewLog("server warmed up").Print()
	}
}

func (a *Application) pullImage() {
	pullOptions := types.ImagePullOptions{}

	if a.Registry != "" {
		pullOptions.RegistryAuth = a.GetRegistry().ToBase64()
	}

	events, err := globals.Docker.ImagePull(context.Background(), a.Image, pullOptions)
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

func (a *Application) StopContainer(name string) bool {
	c := a.getContainer(name)

	if c == nil {
		return false
	}

	err := globals.Docker.ContainerStop(context.Background(), c.ID, nil)
	ui.Check(err)

	err = globals.Docker.ContainerRemove(context.Background(), name, types.ContainerRemoveOptions{})
	ui.Check(err)

	return true
}

func (a *Application) getContainer(name string) *types.Container {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "name", Value: name}),
		All: true,
	})
	ui.Check(err)

	for _, c := range containers {
		// Filters return non-exact matches so {Key: name, Value: example_com_80}
		// could return example_com_80,example_com_80_temporary...qqq
		if c.Names[0] == "/"+name {
			return &c
		}
	}

	return nil
}

func (a *Application) createContainer(name string) string {
	net := a.getNetwork()

	// TODO: Quotas
	// TODO: Better volumes
	ref, err := globals.Docker.ContainerCreate(context.Background(), &container.Config{
		Env:   a.EnvironmentToDockerEnv(),
		Image: a.Image,
		Labels: map[string]string{
			"nest-id":   a.Host + a.Port,
			"nest-host": a.Host,
			"nest-port": a.Port,
		},
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Binds: a.Volumes,
	}, nil, nil, name)
	ui.Check(err)

	err = globals.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
	ui.Check(err)

	// Docker's default network
	_ = globals.Docker.NetworkDisconnect(context.Background(), net.ID, ref.ID, true)
	err = globals.Docker.NetworkConnect(context.Background(), net.ID, ref.ID, nil)
	ui.Check(err)

	return ref.ID
}

func (a *Application) getNetwork() types.NetworkResource {
	networks, err := globals.Docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: a.Network,
		}),
	})
	ui.Check(err)

	// We don't have to check if the network exists, it does, we check for it before we even run the commana.
	net, err := globals.Docker.NetworkInspect(context.Background(), networks[0].ID, types.NetworkInspectOptions{})
	ui.Check(err)
	return net
}

func (a *Application) EnvironmentToDockerEnv() []string {
	variables := make([]string, len(a.Env)-1)

	for k, v := range a.Env {
		variables = append(variables, k+"="+v)
	}

	return variables
}

func (a *Application) waitForContainerToBeHealthy(name string) {
	c := a.getContainer(name)

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

func (a *Application) warmServer(name string) {
	_ = a.getContainer(name).NetworkSettings.Networks[a.Network]

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

func (a Application) ContainerName() string {
	return strings.ReplaceAll(a.Host, ".", "_") + "_" + a.Port
}

func (a Application) TemporaryContainerName() string {
	return a.ContainerName() + "_temporary"
}
