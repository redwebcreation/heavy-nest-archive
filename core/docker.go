package core

import (
	"github.com/docker/docker/client"
	"log"
)

func GetDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		log.Fatal(err)
	}

	return cli
}