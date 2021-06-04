package apply

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func run(cmd *cobra.Command, _ []string) {
	configFile := core.FindConfig(core.ConfigFile())

	currentChecksum, _ := configFile.Checksum()
	previousChecksum, _ := core.GetKey("previous_checksum")
	fmt.Println("Previous config checksum : " + previousChecksum)
	fmt.Println("Current config checksum : " + currentChecksum)

	force, _ := cmd.Flags().GetBool("force")

	if currentChecksum == previousChecksum {
		if force {
			fmt.Println("No changes. Not aborting as --force is set to true.")
		} else {
			fmt.Println("No changes. Aborting.")
			return
		}
	} else {
		fmt.Println("Found changes.")
	}

	config, _ := configFile.Resolve()

	if len(config.Applications) == 0 {
		fmt.Println("No applications found.")
	}

	for _, application := range config.Applications {
		for _, binding := range application.Bindings {
			fmt.Println("\n[" + binding.Host + "]")

			fmt.Println("  - Cleaning up old state.")
			//Here tell the reverse proxy to use a cloned version of the container
			//for 0 downtime deployment
			err := application.Start(binding, true)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = application.CleanUp(func(container types.Container) bool {
				return !strings.HasSuffix(container.Names[0], "_temporary")
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("  - Old state cleaned up.")
			fmt.Println("  - Starting the containers.")
			err = application.Start(binding, false)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = application.CleanUp(func(container types.Container) bool {
				return strings.HasSuffix(container.Names[0], "_temporary")
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("  - Containers started.")
		}
	}

	core.SetKey("previous_checksum", currentChecksum)
}

func NewCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Applies your configuration to the server.",
		Run:   run,
	}

	applyCmd.Flags().BoolP("force", "f", false, "Force the apply command to run")

	return applyCmd
}
