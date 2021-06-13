package main

import (
	"github.com/redwebcreation/hez2/ansi"
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
		ansi.Text("Hez is not meant to be used on ["+runtime.GOOS+"].", ansi.Orange)
		ansi.Text("Please run the command again with the flag --i-will-break-my-computer.", ansi.Orange)
		return
	}

	_, err := os.Stat(globals.ConfigFile)

	if os.IsNotExist(err) {
		ansi.Text("Configuration file not found at "+globals.ConfigFile, ansi.Red)
		return
	}

	cli.Execute()
}
