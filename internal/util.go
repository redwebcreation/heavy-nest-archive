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

func formatTimeInMilliseconds(t time.Time) string {
	diff := strconv.Itoa(t.Nanosecond() / (int(time.Millisecond) / 100))
	ms := diff[0:len(diff)-2]
	if ms == "" {
		ms = "0"
	}

	precision := diff[len(diff) - 2:]
	return ms + "." + precision + "ms"
}