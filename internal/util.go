package internal

import (
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

type CommandConfigurationHandler func(*cobra.Command)
type CommandHandler func(*cobra.Command, []string) error

func CreateCommand(command *cobra.Command, configurationHandler CommandConfigurationHandler, commandHandler CommandHandler) *cobra.Command {
	if configurationHandler != nil {
		configurationHandler(command)
	}

	command.RunE = commandHandler
	command.SilenceErrors = true
	command.SilenceUsage = true
	return command
}
