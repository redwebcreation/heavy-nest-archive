package core

import (
	"os"
)

func Home() string {
	home, err := os.UserHomeDir()

	if err != nil {
		home, err = os.Getwd()

		return home
	}

	return home
}

func ConfigFile() string {
	return Home() + "/.hez.yml"
}

func configDirectory() string {
	return Home() + "/.hez"
}
