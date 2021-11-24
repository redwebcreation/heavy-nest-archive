package main

import (
	"github.com/spf13/cobra"
	"github.com/wormable/nest/ansi"
	"github.com/wormable/nest/cmd"
)

func main() {
	cli := &cobra.Command{
		Use:   "nest",
		Short: "nest makes orchestrating containers easy.",
		Long:  "nest is to tool to orchestrate containers.",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
		},
	}

	cli.AddCommand(
		cmd.ApplyCommand(),
		cmd.DiagnoseCommand(),
		cmd.SelfUpdateCommand(),
		cmd.ProxyCommand(),
		cmd.InitCommand(),
		cmd.MachineCommand(),
	)

	cli.PersistentFlags().BoolP("no-ansi", "A", false, "Disable ANSI output")

	err := cli.Execute()
	ansi.Check(err)
}
