package main

import (
	"fmt"

	"github.com/redwebcreation/nest/client"
	"github.com/redwebcreation/nest/cmd"
	"github.com/redwebcreation/nest/ui"
	"github.com/spf13/cobra"
)

func main() {
	cli := &cobra.Command{
		Use:   "nest",
		Short: "nest makes orchestrating containers easy.",
		Long:  "nest is to tool to orchestrate containers and manage the environment around them.",
	}

	cli.AddCommand(
		wrap(cmd.ApplyCommand()),
		wrap(cmd.DiagnoseCommand()),
		wrap(cmd.StopCommand()),
	)
	cli.SilenceErrors = true
	err := cli.Execute()
	ui.Check(err)
}

func wrap(cmd *cobra.Command) *cobra.Command {
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "diagnose" {
			return nil
		}

		diagnosis := client.Analyse(client.Config)

		if diagnosis.ErrorCount > 0 {
			return fmt.Errorf("your configuration file is invalid, please run `nest diagnose` for details")
		}

		return nil
	}

	return cmd
}
