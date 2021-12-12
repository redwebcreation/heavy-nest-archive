package config

import (
	"os"

	"github.com/redwebcreation/nest/service"
)

type ServiceConfig struct {
	service.Service
	Hosts   []string `json:"hosts" yaml:"hosts"`
	Include string   `json:"include" yaml:"include"`
}

type Registry struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Resolver interface {
	ReadFile(path string) ([]byte, error)
}

type Config struct {
	resolver   Resolver
	Services   []ServiceConfig `json:"services" yaml:"services"`
	Registries []struct {
		Registry
		Include string `json:"include" yaml:"include"`
	}
	Vault struct {
		// Driver Must be one of the following: "file"
		//
		Driver string `json:"driver" yaml:"driver"`
		// Options is a map of driver specific options
		// For the file driver, the following options are available:
		// - path: the path to the vault directory
		Options interface{} `json:"options" yaml:"options"`
	} `json:"vault" yaml:"vault"`
}

type FsResolver struct{}

func (f FsResolver) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// nest apply [--from=$(pwd)/config.yml] [--from-remote=<git repository>/<path>] [--only=<service>] [--exclude=<service>] [--dry-run] [--remote-branch=<branch>] [--remote-commit=<commit>]++

// nest new [registry|service] [path=] [--dry-run] [--create-last]
// Interactively create a new service or registry
// The --create-last, if set, will create the service or registry that was created last.
// Find a better name than create-last.
// Useful if you dry run first (you should) and don't want to type everything again.

// nest self-update [version=latest] [--dry-run]

// nest check [--from=$(pwd)/config.yml] [--from-remote=<git repository>/<path>]
// former doctor
// + vault driver is valid
// + vault driver options are valid for the given driver
// + make sure that include statements point to a file. n
// + check if all the secrets required exist in the vault.
// + check if all environments are valid / exist.

// nest vault put [name] [--force]
// stdin: [password]
// stdin: [value]

// nest vault get [name]
// stdin: [password]
// stdout: [value]

// nest vault delete [name]
