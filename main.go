package main

import (
	"github.com/redwebcreation/hez/cli/apply"
	"github.com/redwebcreation/hez/cli/proxy"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	hezCli := &cobra.Command{
		Use:   "hez",
		Short: "Hez makes orchestrating containers easy.",
		Long:  `Hez is a tool to orchestrate containers and manage the environment around it.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			core.EnsureConfigIsValid()
		},
	}

	hezCli.AddCommand(apply.NewCommand())
	hezCli.AddCommand(proxy.NewCommand())

	if err := hezCli.Execute(); err != nil {
		os.Exit(1)
	}
}
