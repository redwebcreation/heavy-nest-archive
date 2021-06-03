package cli

import (
	"github.com/redwebcreation/hez/cli/apply"
	"github.com/redwebcreation/hez/cli/config"
	"github.com/redwebcreation/hez/cli/env"
	"github.com/redwebcreation/hez/cli/proxy"
	"github.com/redwebcreation/hez/cli/stop"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
)

var rootCli = &cobra.Command{
	Use:   "hez",
	Short: "Hez makes orchestrating containers easy.",
	Long:  `Hez is a tool to orchestrate containers and manage the environment around it.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "new" {
			core.EnsureConfigIsValid()
		}
	},
}

func Execute() {
	rootCli.AddCommand(env.NewCommand())
	rootCli.AddCommand(proxy.NewCommand())
	rootCli.AddCommand(config.NewCommand())
	rootCli.AddCommand(apply.NewCommand())
	rootCli.AddCommand(stop.NewCommand())
	cobra.CheckErr(rootCli.Execute())
}
