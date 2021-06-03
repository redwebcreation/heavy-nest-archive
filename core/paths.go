package core

import (
	"os"
)

func ConfigFile() string {
	return ConfigDirectory() + "/hez.yml"
}

func ConfigDirectory() string {
	return "/etc/hez"
}

func StorageDirectory() string {
	home, _ := os.UserHomeDir()

	return home + "/.config/hez/storage"
}
