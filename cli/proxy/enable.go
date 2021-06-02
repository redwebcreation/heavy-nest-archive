package proxy

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func runEnableCommand(_ *cobra.Command, _ []string) {
	config, _ := core.GetConfig()

	if !core.IsRunningAsRoot() {
		fmt.Println("This command requires elevated privileges.")
		os.Exit(1)
	}

	isProxyEnabled, err := core.IsProxyEnabled()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if isProxyEnabled {
		fmt.Println("Proxy is already enabled.")
		os.Exit(1)
	}

	err = core.EnableProxy(config.Proxy)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Proxy has been successfully enabled.")
}

func initEnableCommand() *cobra.Command {
	enableCommand := &cobra.Command{
		Use:   "enable",
		Short: "Enables the reverse proxy.",
		Run:   runEnableCommand,
	}

	return enableCommand
}
