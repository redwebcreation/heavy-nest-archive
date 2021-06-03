package core

import (
	"os"
	"strings"
)

func EnvironmentPath(path string) string {
	if path == "" {
		return ConfigDirectory() + "/environments"
	}

	return ConfigDirectory() + "/environments/" + strings.TrimLeft(path, "/")
}

func EnvironmentExists(environment string) bool {
	path := EnvironmentPath(environment)
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}
