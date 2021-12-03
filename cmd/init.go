package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/ansi"
	"github.com/wormable/nest/common"
	"github.com/wormable/nest/globals"
	"os"
)

var force bool

func runInitCommand(_ *cobra.Command, _ []string) error {
	_, err := os.Stat(common.ConfigFile)

	if !force && os.IsNotExist(err) == false {
		return fmt.Errorf("config file already exists")
	}

	// create new file
	err = os.WriteFile(common.ConfigFile, globals.DefaultConfig, 0644)
	if err != nil {
		return err
	}

	fmt.Printf(ansi.Green.Fg()+"%s created (%d bytes written)\n"+ansi.Reset, common.ConfigFile, len(globals.DefaultConfig))

	return nil
}

func InitCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "init",
		Short: "create a new configuration file",
	}, runInitCommand, func(command *cobra.Command) {
		command.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing configuration file")
	})
}
