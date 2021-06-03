package config

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

func runNewCommand(_ *cobra.Command, _ []string) {
	errors := core.IsConfigValid()

	if len(errors) == 0 {
		fmt.Println("The configuration has already been created and contains no errors.")
		os.Exit(1)
	}

	configFileExists, configDirectoryExists := true, true

	for _, err := range errors {
		if strings.Contains(err.Error(), "stat /etc/hez/hez.yml") {
			configFileExists = false
		} else if strings.Contains(err.Error(), "stat /etc/hez:") {
			configDirectoryExists = false
		}
	}

	if !(!configFileExists && !configDirectoryExists) || (configFileExists && configDirectoryExists && len(errors) > 0) {
		fmt.Print("It seems like your configuration is in a partially correct state.\n\n")

		fmt.Println("Config file exists : " + strconv.FormatBool(configFileExists))
		fmt.Println("Config directory exists : " + strconv.FormatBool(configDirectoryExists))
		fmt.Println("Total errors : " + strconv.FormatInt(int64(len(errors)), 10) + "\n")

		fmt.Println("You may run `hez diagnose` to get more details.")

		fmt.Println("\nNote: The `new` command is made to be run once after installing Hez.")
		os.Exit(1)
	}

	if !core.IsRunningAsRoot() {
		fmt.Println("This command requires elevated privileges.")
		os.Exit(1)
	}

	toCreate := [3]string{
		core.ConfigDirectory(),
		core.ConfigDirectory() + "/environments",
		core.ConfigDirectory() + "/ssl",
	}

	for _, directoryPath := range toCreate {
		err := os.Mkdir(directoryPath, os.FileMode(0777))

		if err != nil {
			fmt.Println("For [" + directoryPath + "] :" + err.Error())
			Rollback()
			os.Exit(1)
		}

		fmt.Println("Creating [" + directoryPath + "].")
	}

	err := os.WriteFile(core.ConfigFile(), []byte(getDefaultConfigContents()), os.FileMode(0777))

	if err != nil {
		fmt.Println(err)
		Rollback()
		os.Exit(1)
	}

	fmt.Println("Creating [" + core.ConfigFile() + "]")
	fmt.Println("Configuration successfully created.")
}

func getDefaultConfigContents() string {
	return strings.TrimSpace(`
logs:
  level: 0
  redirections:
    - stdout
proxy:
  port: 80
  ssl: 443
  self_signed: false
applications: []
`)
}

func initRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Creates an initial configuration",
		Run:   runNewCommand,
	}
}

func Rollback() {
	_ = os.RemoveAll(core.ConfigDirectory())
	fmt.Println("Rolling back. ")
}
