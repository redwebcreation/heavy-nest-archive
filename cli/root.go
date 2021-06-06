package cli

import (
	"fmt"
	"github.com/redwebcreation/hez/cli/apply"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

var rootCli = &cobra.Command{
	Use:   "hez",
	Short: "Hez makes orchestrating containers easy.",
	Long:  `Hez is a tool to orchestrate containers and manage the environment around it.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "new" {
			valid := core.FindConfig(core.ConfigFile()).IsValid()

			if !valid {
				fmt.Println("Invalid configuration.")
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	//rootCli.AddCommand(proxy.NewCommand())
	//rootCli.AddCommand(config.NewCommand())
	rootCli.AddCommand(apply.NewCommand())
	//rootCli.AddCommand(stop.NewCommand())
	//rootCli.AddCommand(health.NewCommand())
	cobra.CheckErr(rootCli.Execute())
}
