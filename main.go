package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/ansi"
	"github.com/wormable/nest/cmd"
	"github.com/wormable/nest/cmd/api/tokens"
	"github.com/wormable/nest/common"
	"os"
)

func main() {
	isInitCommand := false

	for _, arg := range os.Args {
		if arg == "init" {
			isInitCommand = true
			break
		}
	}

	_, err := os.Stat(common.ConfigFile)

	if os.IsNotExist(err) && !isInitCommand {
		fmt.Printf("%sConfiguration file not found at %s.%s\n", ansi.Red.Fg(), common.ConfigFile, ansi.Reset)
		fmt.Printf("%sYou can create one by running `nest init` with elevated privileges.%s\n", ansi.Red.Fg(), ansi.Reset)
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		common.LoadConfig()
	}

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

	api := &cobra.Command{
		Use:   "api",
		Short: "manage the integrated deployment api",
	}

	api.AddCommand(
		tokens.RootCommand(),
	)

	cli.AddCommand(
		cmd.ApplyCommand(),
		cmd.DiagnoseCommand(),
		cmd.SelfUpdateCommand(),
		cmd.ProxyCommand(),
		cmd.InitCommand(),
		cmd.MachineCommand(),
		api,
	)

	cli.PersistentFlags().Bool("no-ansi", false, "Disable ANSI output")

	cli.SetHelpCommand(&cobra.Command{
		Use:    "_help_",
		Hidden: true,
	})

	err = cli.Execute()
	ansi.Check(err)
}
