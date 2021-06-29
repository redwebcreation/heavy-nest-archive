package main

import (
	"github.com/redwebcreation/hez/cli"
	"github.com/redwebcreation/hez/internal"
	ansi2 "github.com/redwebcreation/hez/internal/ui"
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
		ansi2.Warning("Hez is not meant to be used on [" + runtime.GOOS + "].")
		ansi2.Warning("Please run the command again with the flag --i-will-break-my-computer.")
		os.Exit(1)
	}

	_, err := os.Stat(internal.ConfigFile)

	if os.IsNotExist(err) {
		ansi2.Error("Configuration file not found at " + internal.ConfigFile)
		os.Exit(1)
	}

	cli.Execute()
}
