package config

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func runNewCommand(_ *cobra.Command, _ []string) {
	isConfigValid := core.FindConfig(core.ConfigFile()).IsValid()

	if isConfigValid {
		fmt.Println("The configuration already exists and seems to be valid.")
		os.Exit(1)
	}

	if !core.IsRunningAsRoot() {
		fmt.Println("This command requires elevated privileges.")
		os.Exit(1)
	}

	_ = os.MkdirAll(core.ConfigDirectory(), os.FileMode(0777))
	err := os.WriteFile(core.ConfigFile(), []byte(getDefaultConfigContents()), os.FileMode(0777))

	if err != nil {
		fmt.Println(err)
		_ = os.RemoveAll(core.ConfigDirectory())
		fmt.Println("Rolling back.")
		os.Exit(1)
	}

	fmt.Println("Creating [" + core.ConfigFile() + "]")
	fmt.Println("Configuration successfully created.")
}

func getDefaultConfigContents() string {
	return strings.TrimSpace(`
network: bridge
applications: []
proxy:
  port: 80
  ssl: 443
  self_signed: false
  logs:
    level: 0
    redirections:
      - stdout
`)
}

func initRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Creates an initial configuration",
		Run:   runNewCommand,
	}
}
