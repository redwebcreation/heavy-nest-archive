package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies your configuration to the server.",
	Run: func(cmd *cobra.Command, args []string) {
		currentChecksum := core.GetConfigChecksum()
		previousChecksum := core.GetKey("previous_checksum")
		fmt.Printf("Previous config checksum : %x\n", previousChecksum)
		fmt.Printf("Current config checksum : %x\n", currentChecksum)

		force, _ := cmd.Flags().GetBool("force")

		if currentChecksum == previousChecksum {
			if force {
				fmt.Println("No changes to the config. Not aborting as --force is set to true.")
			} else {
				fmt.Println("No changes. Aborting.")
				return
			}
		}

		conf, err := core.GetConfig()

		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		for _, application := range conf.Applications {
			fmt.Println("\n[" + application.Name + "]")

			wasRunning := application.HasRunningContainer()

			if wasRunning {
				fmt.Println("  - Found a running container.")

				stoppedContainer := application.GetContainer()

				application.StopContainer()

				fmt.Println("  - Stopped the running container.")

				err := core.GetDockerClient().ContainerRemove(context.Background(), stoppedContainer.ID, types.ContainerRemoveOptions{
					RemoveVolumes: false,
					RemoveLinks:   false,
					Force:         false,
				})

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println("  - Deleting the old container [" + application.Name + "].")
			}

			if wasRunning {
				fmt.Println("  - Restarting the container.")
			} else {
				fmt.Println("  - Starting the container.")
			}
			application.Start()
			fmt.Println("  - Container restarted.")
		}

		core.SetKeyOverride("previous_checksum", currentChecksum)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().BoolP("force", "f", false, "Force the apply command to run")
}
