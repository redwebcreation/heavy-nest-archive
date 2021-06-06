package health

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func run(_ *cobra.Command, _ []string) {
	proxiablesContainers, err := core.GetProxiableContainers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(proxiablesContainers) == 0 {
		fmt.Println("No running containers.")
	}

	for _, proxiable := range proxiablesContainers {
		var healthStatus string

		if proxiable.Container.State.Health == nil {
			healthStatus = "healthy"
		} else {
			healthStatus = proxiable.Container.State.Health.Status
		}

		if healthStatus == "healthy" {
			fmt.Println(proxiable.Name + ": " + healthStatus)
		} else {
			fmt.Println("[" + proxiable.Name + "]")
			fmt.Printf("failing_streak=%d\n", proxiable.Container.State.Health.FailingStreak)

			lastLog := proxiable.Container.State.Health.Log[len(proxiable.Container.State.Health.Log)-1]
			fmt.Println("started_at=" + lastLog.Start.String())
			fmt.Println("ended_at=" + lastLog.End.String())
			fmt.Printf("exit_code=%d\n", lastLog.ExitCode)
		}
	}
}

func NewCommand() *cobra.Command {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Returns the containers' health",
		Run:   run,
	}

	return healthCmd

}
