package core

import (
	"fmt"
	"github.com/docker/docker/client"
	"os"
)

func GetDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return cli
}
