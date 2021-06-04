package proxy

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func runEnableCommand(_ *cobra.Command, _ []string) {
	config, _ := core.FindConfig(core.ConfigFile()).Resolve()

	if !core.IsRunningAsRoot() {
		fmt.Println("This command requires elevated privileges.")
		os.Exit(1)
	}

	if core.IsProxyEnabled() {
		fmt.Println("Proxy is already enabled.")
		os.Exit(1)
	}

	if !*config.Proxy.SelfSigned {
		config.Proxy.SelfSigned = &selfSigned
	}

	err := core.EnableProxy(config.Proxy)

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

	enableCommand.Flags().BoolVar(&selfSigned, "self-signed", false, "Force the use of self signed certificates.")

	return enableCommand
}
