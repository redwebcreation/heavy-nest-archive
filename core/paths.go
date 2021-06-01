package core

import (
	"fmt"
	"os"
)

func Home() string {
	home, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	return home
}

func ConfigFile() string {
	configFile := Home() + "/.hez.yml"

	_, err := os.Stat(configFile)

	if os.IsNotExist(err) {
		fmt.Println("The config file " + configFile + " does not exist.")
		os.Exit(1)
	}

	return configFile
}

func ConfigDirectory() string {
	configDirectory := Home() + "/.hez"

	EnsureDirectoryExists(configDirectory)

	return configDirectory
}
