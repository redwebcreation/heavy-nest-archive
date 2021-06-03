package config

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func runDeleteCommand(_ *cobra.Command, _ []string) {
	config := core.GetConfig()

	if !core.IsRunningAsRoot() {
		fmt.Println("This command requires elevated privileges.")
		os.Exit(1)
	}

	fmt.Println("[proxy]")
	if core.IsProxyEnabled() {
		fmt.Println("  - The proxy is running.")
		fmt.Println("  - You won't be able to disable it once you delete the configuration.")
		fmt.Println("  - Please stop the reverse proxy first.")
		os.Exit(1)
	} else {
		fmt.Println("  - The proxy was already disabled.")
	}

	fmt.Println("[container]")
	if len(config.Applications) == 0 {
		fmt.Println("  - No applications found.")
	} else {
		for _, application := range config.Applications {
			err := application.CleanUp()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	fmt.Println("[config]")
	_ = os.RemoveAll(core.ConfigDirectory())
	fmt.Println("  - Deleting [" + core.ConfigDirectory() + "]")
	fmt.Println("The configuration has been successfully deleted.")
}

func initDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Removes all the configuration.",
		Run:   runDeleteCommand,
	}
}
