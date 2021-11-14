package cmd

import "github.com/spf13/cobra"

func runInitCommand(_ *cobra.Command, _ []string) error {
	return nil
}

func InitCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuraiton file",
	}, nil, runInitCommand)
}
