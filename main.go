package main

import (
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
	"github.com/wormable/nest/cmd/ui"
)

func main() {
	cli := &cobra.Command{
		Use:   "nest",
		Short: "nest makes orchestrating containers easy.",
		Long:  "nest is to tool to orchestrate containers and manage the environment around them.",
	}

	// TODO: Move wrap to create command
	cli.AddCommand(
		cmd.ApplyCommand(),
		cmd.DiagnoseCommand(),
		cmd.StopCommand(),
		cmd.SelfUpdateCommand(),
	)

	err := cli.Execute()
	ui.Check(err)
}
