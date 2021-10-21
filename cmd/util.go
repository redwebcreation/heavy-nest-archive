package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/common"
)

type CommandConfigurationHandler func(*cobra.Command)
type CommandHandler func(*cobra.Command, []string) error

func CreateCommand(command *cobra.Command, configurationHandler CommandConfigurationHandler, commandHandler CommandHandler) *cobra.Command {
	if configurationHandler != nil {
		configurationHandler(command)
	}

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "diagnose" {
			return nil
		}

		diagnosis := common.Analyse(common.Config)

		if diagnosis.ErrorCount > 0 {
			return fmt.Errorf("your configuration file is invalid, please run `nest diagnose` for details")
		}

		return nil
	}
	command.RunE = commandHandler
	command.SilenceErrors = true
	command.SilenceUsage = true
	return command
}

func ElevateProcess() {
	cmd := exec.Command("sudo", "ls")
	cmd.Stdin = os.Stdin
	cmd.Run()
}
