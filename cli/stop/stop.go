package stop

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
)

func run(_ *cobra.Command, _ []string) {
	config, _ := core.FindConfig(core.ConfigFile()).Resolve()

	if len(config.Applications) == 0 {
		fmt.Println("No containers running.")
		return
	}

	for _, application := range config.Applications {
		container, err := application.RemoveApplicationContainer()

		if container == "" && err == nil {
			continue
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Removing [" + application.Name(false) + "]")

		_, _ = application.RemoveEphemeralContainer()
	}

	fmt.Println("Stopped all containers successfully.")
}

func NewCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops all running containers",
		Run:   run,
	}

	return applyCmd
}
