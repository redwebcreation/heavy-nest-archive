package core

import (
	"fmt"
	"os"
)

func ConfigFile() string {
	return ConfigDirectory() + "/hez.yml"
}

func ConfigDirectory() string {
	return "/etc/hez"
}

func StorageDirectory() string {
	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home + "/.config/hez/storage"
}
