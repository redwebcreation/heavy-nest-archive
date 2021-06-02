package stop

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func run(_ *cobra.Command, _ []string) {
	conf, _ := core.GetConfig()

	if len(conf.Applications) == 0 {
		fmt.Println("No containers running.")
		return
	}

	for _, application := range conf.Applications {
		fmt.Println("[" + application.Name + "]")
		err := application.CleanUp()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
