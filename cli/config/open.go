package config

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func runOpenCommand(_ *cobra.Command, _ []string) {
	if !core.IsRunningAsRoot() {
		fmt.Println("This command requires elevated privileges.")
		os.Exit(1)
	}

	cmd := exec.Command(os.Getenv("EDITOR"), core.ConfigFile())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Opens your config file in your default editor",
		Run:   runOpenCommand,
	}
}
