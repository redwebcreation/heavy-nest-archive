package cli

import (
	"github.com/redwebcreation/hez/ui"
	"github.com/spf13/cobra"
)

func Execute() {
	cli := &cobra.Command{
		Use:   "hez",
		Short: "Hez makes orchestrating containers easy.",
		Long:  "Hez is to tool to orchestrate containers and manage the environment around them.",
	}

	cli.AddCommand(ApplyCommand())

	cli.SilenceErrors = true
	err := cli.Execute()
	ui.Check(err)
}
