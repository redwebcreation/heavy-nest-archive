package proxy

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func runDisableCommand(cmd *cobra.Command, _ []string) {
	config := core.GetConfig()

	if !core.IsRunningAsRoot() {
		fmt.Println("This command requires elevated privileges.")
		os.Exit(1)
	}

	if !core.IsProxyEnabled() {
		fmt.Println("Proxy is already disabled.")
		os.Exit(1)
	}

	err := core.DisableProxy(config.Proxy)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Proxy has been successfully disabled.")
}

func initDisableCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disables the reverse proxy.",
		Run:   runDisableCommand,
	}
}
