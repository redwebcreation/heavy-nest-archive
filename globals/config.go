package globals

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"strings"
)

type HezConfig struct {
	DefaultNetwork string        `yaml:"default_network"`
	Applications   []Application `yaml:"applications"`
	Proxy          struct {
		Logs struct {
			Level        int      `yaml:"level"`
			Redirections []string `yaml:"redirections"`
		} `yaml:"logs"`
		Http struct {
			Enabled bool `yaml:"enabled"`
			Port    int  `yaml:"port"`
		} `yaml:"http"`
		Https struct {
			Enabled    bool  `yaml:"enabled"`
			Port       int   `yaml:"port"`
			SelfSigned *bool `yaml:"self_signed"`
		} `yaml:"https"`
	} `yaml:"proxy"`
}

type Application struct {
	Image         string   `yaml:"image"`
	Host          string   `yaml:"host"`
	ContainerPort int      `yaml:"container_port"`
	Network       string   `yaml:"network"`
	Warm          *bool    `yaml:"warm"`
	Env           []string `yaml:"env"`
	Volumes       []struct {
		From string
		To   string
	} `yaml:"volumes"`
	//Replicas      int
	Cpu      string `yaml:"cpu"`
	Ram      string `yaml:"ram"`
	Firewall struct {
		Rate int `yaml:"rate"` // max requests every 5 minutes per ip.
	} `yaml:"firewall"`
	Hooks struct {
		Before []string `yaml:"before"`
		After  []string `yaml:"after"`
	} `yaml:"hooks"`
}

var Config *HezConfig
var ConfigFile = "/etc/hez/hez.yml"

func init() {
	data := HezConfig{}
	bytes, _ := os.ReadFile(ConfigFile)

	err := yaml.Unmarshal(bytes, &data)

	if err != nil {
		panic(err)
	}

	useDefaults(&data)

	Config = &data
}

func useDefaults(config *HezConfig) {
	if config.DefaultNetwork == "" {
		config.DefaultNetwork = "bridge"
	}

	if config.Proxy.Http.Port == 0 {
		config.Proxy.Http.Port = 80
	}

	if config.Proxy.Https.Port == 0 {
		config.Proxy.Http.Port = 443
	}

	if config.Proxy.Https.SelfSigned == nil {
		selfSigned := false
		config.Proxy.Https.SelfSigned = &selfSigned
	}

	for i := range config.Applications {
		if config.Applications[i].ContainerPort == 0 {
			config.Applications[i].ContainerPort = 80
		}

		if config.Applications[i].Network == "" {
			config.Applications[i].Network = config.DefaultNetwork
		}

		if config.Applications[i].Warm == nil {
			warm := true
			config.Applications[i].Warm = &warm
		}
	}
}

func (application Application) Name() string {
	return strings.ReplaceAll(application.Host, ".", "_") + "_" + strconv.Itoa(application.ContainerPort)
}

func (application Application) NameWithSuffix(suffix string) string {
	return application.Name() + "_" + suffix
}

func (application Application) CreateTemporaryContainer() (string, error) {
	return application.createContainer(
		application.NameWithSuffix("temporary"),
	)
}

func (application Application) CreateApplicationContainer() (string, error) {
	return application.createContainer(
		application.Name(),
	)
}

func (application Application) createContainer(name string) (string, error) {
	var mounts []mount.Mount
	env := []string{
		"VIRTUAL_HOST=" + application.Host,
		"VIRTUAL_PORT=" + strconv.Itoa(application.ContainerPort),
	}

	networkDetails, _ := FindNetwork(application.Network)

	for _, volume := range application.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: volume.From,
			Target: volume.To,
		})
	}

	resp, err := Docker.ContainerCreate(context.Background(), &container.Config{
		Env:   ResolveEnvironmentVariables(env, application.Env),
		Image: application.Image,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		Mounts: mounts,
	}, nil, nil, name)

	if err != nil {
		return "", err
	}

	err = Docker.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})

	if err != nil {
		return "", err
	}

	_ = Docker.NetworkDisconnect(context.Background(), networkDetails.ID, resp.ID, true)

	err = Docker.NetworkConnect(context.Background(), networkDetails.ID, resp.ID, nil)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (application Application) StopApplicationContainer() (types.Container, error) {
	return application.stopContainer(application.Name())
}

func (application Application) StopTemporaryContainer() (types.Container, error) {
	return application.stopContainer(
		application.NameWithSuffix("temporary"),
	)
}

func GetContainer(name string) (types.Container, error) {
	containers, err := Docker.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		return types.Container{}, err
	}

	var currentContainer types.Container

	for _, c := range containers {
		if c.Names[0] == "/"+name {
			currentContainer = c
			break
		}
	}

	return currentContainer, nil
}

func (application Application) stopContainer(name string) (types.Container, error) {
	currentContainer, err := GetContainer(name)

	if err != nil {
		return currentContainer, err
	}

	err = Docker.ContainerStop(context.Background(), currentContainer.ID, nil)

	if err != nil {
		return currentContainer, err
	}

	err = Docker.ContainerRemove(context.Background(), currentContainer.Names[0], types.ContainerRemoveOptions{})

	if err != nil {
		return currentContainer, err
	}

	return currentContainer, nil
}

func ResolveEnvironmentVariables(variables []string, env []string) []string {
	for _, envVariable := range env {
		if strings.Contains(envVariable, "=") {
			variables = append(variables, envVariable)
		} else {
			contents, err := os.ReadFile(envVariable)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, envFileVariable := range strings.Split(string(contents), "\n") {
				trimmed := strings.TrimSpace(envFileVariable)

				if trimmed == "" {
					continue
				}

				variables = append(variables, trimmed)
			}
		}
	}

	return variables
}

func FindNetwork(name string) (types.NetworkResource, error) {
	networks, err := Docker.NetworkList(context.Background(), types.NetworkListOptions{})

	if err != nil {
		return types.NetworkResource{}, err
	}

	var network types.NetworkResource

	for _, maybeNetwork := range networks {
		if maybeNetwork.Name == name {
			network = maybeNetwork
			break
		}
	}

	networkDetails, err := Docker.NetworkInspect(context.Background(), network.ID, types.NetworkInspectOptions{})

	if err != nil {
		return networkDetails, err
	}

	return networkDetails, nil
}
