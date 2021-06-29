package globals

import "github.com/docker/docker/client"

var Docker *client.Client

func init() {
	client, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	Docker = client
}
