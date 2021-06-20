package main

import (
	"github.com/redwebcreation/hez/ansi"
	"github.com/redwebcreation/hez/cli"
	"github.com/redwebcreation/hez/core"
	"os"
	"runtime"
)

func main() {
	allowsUntestedOs := false

	for i, arg := range os.Args {
		if arg == "--i-will-break-my-computer" {
			// Removes the argument so Cobra doesn't freak out.
			os.Args = os.Args[:i+copy(os.Args[i:], os.Args[i+1:])]
			allowsUntestedOs = true
			break
		}
	}

	if runtime.GOOS != "linux" && !allowsUntestedOs {
		ansi.Warning("Hez is not meant to be used on [" + runtime.GOOS + "].")
		ansi.Warning("Please run the command again with the flag --i-will-break-my-computer.")
		os.Exit(1)
	}

	_, err := os.Stat(core.ConfigFile)

	if os.IsNotExist(err) {
		ansi.Error("Configuration file not found at " + core.ConfigFile)
		os.Exit(1)
	}

	cli.Execute()
}
