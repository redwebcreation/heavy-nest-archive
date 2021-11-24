package cmd

import "github.com/spf13/cobra"

func runInitCommand(_ *cobra.Command, _ []string) error {
	return nil
}

func InitCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "init",
		Short: "create a new configuration file",
	}, runInitCommand, nil)
}
