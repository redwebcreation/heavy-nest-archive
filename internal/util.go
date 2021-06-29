package internal

import (
	"fmt"
	"github.com/redwebcreation/hez/globals"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func ElevateProcess() error {
	cmd := exec.Command("sudo", "touch", "/tmp/upgrade-process")

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

type CommandConfigurationHandler func(command *cobra.Command)
type CommandHandler func(_ *cobra.Command, _ []string) error

func CreateCommand(command *cobra.Command, commandConfigurationHandler CommandConfigurationHandler, Handler CommandHandler) *cobra.Command {
	if commandConfigurationHandler != nil {
		commandConfigurationHandler(command)
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")

		if showVersion {
			fmt.Println("Hez " + globals.Version)
			return nil
		}

		if Handler == nil {
			return nil
		}

		err := Handler(cmd, args)

		return err
	}
	command.SilenceErrors = true
	command.SilenceUsage = true
	return command
}
