package main

import (
	"github.com/redwebcreation/hez/cli"
	"github.com/redwebcreation/hez/internal"
	"github.com/redwebcreation/hez/internal/ui"
	"os"
	"runtime"
)

func main() {
	allowsUntestedOs := false
	isInitCommand := false

	for i, arg := range os.Args {
		if arg == "--it-wont-work" {
			// Removes the argument so Cobra doesn't freak out.
			os.Args = os.Args[:i+copy(os.Args[i:], os.Args[i+1:])]
			allowsUntestedOs = true
			break
		}

		if i == 1 && arg == "init" {
			isInitCommand = true
		}
	}

	if runtime.GOOS != "linux" && !allowsUntestedOs {
		ui.Warning("Hez is not meant to be used on [" + runtime.GOOS + "].")
		ui.Warning("Please run the command again with the flag--it-wont-work.")
		os.Exit(1)
	}

	_, err := os.Stat(internal.ConfigFile)

	if os.IsNotExist(err) && !isInitCommand {
		ui.Error("Configuration file not found at " + internal.ConfigFile)
		ui.Error("You can create one with `hez init` in an terminal with elevated privileges.")
		os.Exit(1)
	}

	cli.Execute()
}
