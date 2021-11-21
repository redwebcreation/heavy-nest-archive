package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/wormable/nest/ansi"
	"github.com/wormable/nest/globals"
)

type Application struct {
	Image     string            `json:"image,omitempty"`
	Host      string            `json:"host,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	EnvFiles  []string          `json:"env_files,omitempty"`
	Volumes   []string          `json:"volumes,omitempty"`
	Aliases   []string          `json:"aliases,omitempty"`
	Warm      bool              `json:"warm,omitempty"`
	LogPolicy string            `json:"log_policy,omitempty"`
	Registry  string            `json:"registry,omitempty"`
	Network   string            `json:"network,omitempty"`
	Port      string            `json:"port,omitempty"`
	Hooks     struct {
		PreStart []string `json:"pre_start,omitempty"`
	} `json:"hooks,omitempty"`
}

func (a Application) GetRegistry() *RegistryConfiguration {
	if a.Registry == "" {
		return nil
	}

	for _, registry := range Config.Registries {
		if registry.Name == a.Registry {
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
		ansi.NewLog("pulling %s", a.Image).Render()
		a.pullImage()
		ansi.NewLog("successfully downloaded %s", a.Image).Render()
	}

	ansi.NewLog("stopping container %s", opts.Name).Render()
	stopped := a.StopContainer(opts.Name)

	if stopped {
		ansi.CursorUp(1)
		ansi.ClearLine()
		ansi.NewLog("stopped container %s", opts.Name).Render()
	} else {
		ansi.CursorUp(3)
		ansi.ClearLine()
	}

	ansi.NewLog("creating container %s", opts.Name).Render()
	id := a.createContainer(opts.Name)
	if len(a.Hooks.PreStart) > 0 {
		a.RunPostStartHook(id)
		ansi.NewLog("created container %s", opts.Name).Render()
	} else {
		ansi.CursorUp(1)
		ansi.ClearLine()
		ansi.NewLog("created container %s", opts.Name).Render()
	}

	if opts.Healthchecks {
		a.waitForContainerToBeHealthy(opts.Name)
	} else {
		ansi.NewLog("skipping healthchecks").SetColor(ansi.Gray).SetPrefix(ansi.Gray.Fg() + " - " + ansi.Reset).Render()
	}

	if a.Warm {
		ansi.NewLog("warming up server").Render()
		a.warmServer(opts.Name)
		ansi.NewLog("server warmed up").Render()
	}
}

func (a *Application) pullImage() {
	pullOptions := types.ImagePullOptions{}

	if a.Registry != "" {
		pullOptions.RegistryAuth = a.GetRegistry().ToBase64()
	}

	events, err := globals.Docker.ImagePull(context.Background(), a.Image, pullOptions)
	ansi.Check(err)

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

	progress := ansi.NewProgress(60, 50).SetPrefix("    " + "    " + ansi.White.Fg() + "[")

	fmt.Println()
	progress.Render()

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			ansi.Check(err)
		}

		progress.WithLabel(ansi.Gray.Fg() + strings.Replace(strings.ToLower(event.Status), "status: ", "", 1) + ansi.Reset).Increment()
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
	ansi.Check(err)

	err = globals.Docker.ContainerRemove(context.Background(), name, types.ContainerRemoveOptions{})
	ansi.Check(err)

	return true
}

func (a Application) getRunningContainer() *types.Container {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "name",
				Value: a.ContainerName(), // returns all containers starting with ContainerName() including TemporaryContainerName()
			}),
	})
	ansi.Check(err)

	for _, c := range containers {
		if c.Names[0] == "/"+a.ContainerName() || c.Names[0] == "/"+a.TemporaryContainerName() {
			return &c
		}
	}

	return nil
}

func (a *Application) getContainer(name string) *types.Container {
	containers, err := globals.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "name", Value: name}),
		All: true,
	})
	ansi.Check(err)

	for _, c := range containers {
		// Filters return non-exact matches so {Key: name, Value: example_com_80}
		// could return example_com_80,example_com_80_temporary...
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
	ansi.Check(err)

	err = globals.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
	ansi.Check(err)

	_ = globals.Docker.NetworkDisconnect(context.Background(), net.ID, ref.ID, true)
	err = globals.Docker.NetworkConnect(context.Background(), net.ID, ref.ID, nil)
	ansi.Check(err)

	return ref.ID
}

func (a *Application) RunPostStartHook(id string) {
	ansi.NewLog("running post_start hooks").Render()

	for _, hook := range a.Hooks.PreStart {
		ansi.NewLog("%s", hook).SetColor(ansi.White).SetPrefix("  - ").Render()
		out, err := a.executeHook(hook, id)

		lines := strings.Split(string(out), "\n")

		for _, line := range lines {
			if line == "" {
				continue
			}

			fmt.Printf("%s       | %s\n", ansi.Gray.Fg(), line+ansi.Reset)
		}
		fmt.Println()

		if err != nil {
			ansi.Check(fmt.Errorf("error running %s", hook))
		}
	}
}

func (a *Application) getNetwork() types.NetworkResource {
	networks, err := globals.Docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: a.Network,
		}),
	})
	ansi.Check(err)

	// We don't have to check if the network exists, it does, we check for it before we even run the commana.
	net, err := globals.Docker.NetworkInspect(context.Background(), networks[0].ID, types.NetworkInspectOptions{})
	ansi.Check(err)
	return net
}

func (a *Application) EnvironmentToDockerEnv() []string {
	if len(a.Env) == 0 {
		return []string{}
	}

	variables := make([]string, len(a.Env)-1)

	for k, v := range a.Env {
		variables = append(variables, k+"="+v)
	}

	return variables
}

func (a *Application) waitForContainerToBeHealthy(name string) {
	c := a.getContainer(name)

	inspection, err := globals.Docker.ContainerInspect(context.Background(), c.ID)
	ansi.Check(err)

	if inspection.Config.Healthcheck == nil {
		return
	}

	ansi.NewLog("checking the container healthiness").Render()

	maxWaitingTime := inspection.Config.Healthcheck.Interval.Seconds()*math.Max(1.0, float64(inspection.Config.Healthcheck.Retries)) + inspection.Config.Healthcheck.Timeout.Seconds()*math.Max(1.0, float64(inspection.Config.Healthcheck.Retries))

	fmt.Printf(ansi.Gray.Fg()+"      max waiting time: %ds\n\n"+ansi.Reset, int(maxWaitingTime))

	progress := ansi.NewProgress(int(maxWaitingTime), 50)

	for i := maxWaitingTime; i >= 0; i-- {
		if inspection.State.Health.Status == "healthy" {
			break
		}

		time.Sleep(time.Second)
		progress.WithLabel(
			fmt.Sprintf("state: %s, failing streak: %d", inspection.State.Health.Status, inspection.State.Health.FailingStreak),
		).Increment()
		inspection, _ = globals.Docker.ContainerInspect(context.Background(), c.ID)
	}

	progress.Finish()

	// TODO:
	if inspection.State.Health.Status == "healthy" {
		fmt.Println("healthy")
		return
	}

	fmt.Println("UNHEALTHY")
}

func (a *Application) Ip() (string, error) {
	c := a.getRunningContainer()

	if c == nil {
		return "", fmt.Errorf("dead proxy")
	}

	return c.NetworkSettings.Networks[a.Network].IPAddress + ":" + a.Port, nil
}

func (a *Application) Url() (*url.URL, error) {
	ip, err := a.Ip()

	if err != nil {
		return nil, err
	}

	return url.Parse("http://" + ip)
}

func (a *Application) warmServer(name string) {
	ip := a.getContainer(name).NetworkSettings.Networks[a.Network].IPAddress

	var responseTimes [10]float64

	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := http.Get("http://" + ip + ":" + a.Port)
		end := time.Now()
		ansi.Check(err)

		responseTime := end.Add(time.Duration(-start.Nanosecond()))

		responseTimes[i] = float64(responseTime.Nanosecond()) / float64(time.Millisecond)
		responseTimeString := fmt.Sprintf("%.2fms", responseTimes[i])

		fmt.Printf(
			"      %sGET / in%s %s\n",
			ansi.Gray.Fg(),
			ansi.Reset,
			ansi.Blue.Fg()+responseTimeString+ansi.Reset,
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
	if responseTimes[0]-responseTimes[9] <= 0.02 {
		fmt.Print(ansi.Yellow.Fg())
		fmt.Println("      It seems like your server does not need to warmup.")
		fmt.Println("      Consider removing this option in your config to speed up the process.")
		fmt.Print(ansi.Reset)
	}

	fmt.Printf("      %sGain: %.2fms%s Diff: %.2fms Max: %.2fms Min: %.2fms Avg: %.2fms\n\n", ansi.White.Fg(), diff, ansi.Reset, responseTimes[0]-responseTimes[9], max, *min, sum/10.0)
}

func (a Application) ContainerName() string {
	return strings.ReplaceAll(a.Host, ".", "_") + "_" + a.Port
}

func (a Application) TemporaryContainerName() string {
	return a.ContainerName() + "_temporary"
}

func (a Application) executeHook(command string, container string) ([]byte, error) {
	// execute command in container
	cmd := exec.Command("docker", "exec", container, "/bin/sh", "-c", command)
	return cmd.CombinedOutput()
}
