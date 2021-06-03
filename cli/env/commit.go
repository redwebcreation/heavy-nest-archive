package env

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func runCommitCommand(_ *cobra.Command, _ []string) {
	config := core.GetConfig()
	for _, application := range config.Applications {
		currentChecksum := core.GetChecksumForBytes(application.GetCurrentEnvironment())
		stagingChecksum := core.GetChecksumForBytes(application.GetStagingEnvironment())

		if currentChecksum == stagingChecksum {
			fmt.Println("The environment [" + application.Env + "] is in sync.")
		}

		if currentChecksum != stagingChecksum {
			fmt.Println("The environment [" + application.Env + "] is out of sync.")
			fmt.Println()

			fmt.Println("[" + application.Env + "]")
			fmt.Printf("Current: %x\n", currentChecksum)
			fmt.Printf("Staging: %x\n", stagingChecksum)

			err := os.Truncate(core.EnvironmentPath(application.Env+"/current/.env"), 0)
			if err != nil {
				fmt.Println("Could not truncate the current environment file.")
				fmt.Println(err)
				fmt.Println("Aborting.")
				os.Exit(1)
			}

			err = os.WriteFile(core.EnvironmentPath(application.Env+"/current/.env"), application.GetStagingEnvironment(), os.FileMode(0777))
			if err != nil {
				fmt.Println(err)
				fmt.Println("Could not write the new environment file")
				// TODO: Should rollback to the previous env there
				fmt.Println("Aborting.")
				os.Exit(1)
			}

			fmt.Println("The environment [" + application.Env + "] is now in sync.")
		}
	}
	fmt.Println()
	fmt.Println("Committed.")
}

func initCommitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Commits the staging environments to production.",
		Run:   runCommitCommand,
	}
}
