package main

import (
	"github.com/redwebcreation/hez2/cli"
	"github.com/redwebcreation/hez2/globals"
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
		globals.Ansi.Warning("Hez is not meant to be used on [" + runtime.GOOS + "].")
		globals.Ansi.Warning("Please run the command again with the flag --i-will-break-my-computer.")
		return
	}

	_, err := os.Stat(globals.ConfigFile)

	if os.IsNotExist(err) {
		globals.Ansi.ErrorBlock("Configuration file not found at " + globals.ConfigFile)
		return
	}

	cli.Execute()
}
