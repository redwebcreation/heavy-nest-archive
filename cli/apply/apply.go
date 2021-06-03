package apply

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func run(cmd *cobra.Command, _ []string) {
	currentChecksum := core.GetConfigChecksum()
	previousChecksum := core.GetKey("previous_checksum", true)
	fmt.Printf("Previous config checksum : %x\n", previousChecksum)
	fmt.Printf("Current config checksum : %x\n", currentChecksum)

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

	config := core.GetConfig()

	if len(config.Applications) == 0 {
		fmt.Println("No applications found.")
	}

	for _, application := range config.Applications {
		fmt.Println("\n[" + application.Env + "]")

		fmt.Println("  - Cleaning up old state.")
		// Here tell the reverse proxy to use a cloned version of the container
		// for 0 downtime deployment
		err := application.CleanUp()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("  - Old state cleaned up.")
		fmt.Println("  - Starting the containers.")
		err = application.Start()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("  - Containers started.")
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
