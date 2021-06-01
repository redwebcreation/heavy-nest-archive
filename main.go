package main

import (
	"github.com/redwebcreation/hez/cli/apply"
	"github.com/redwebcreation/hez/cli/proxy"
	"github.com/redwebcreation/hez/cli/ssl"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	hezCli := &cobra.Command{
		Use:   "hez",
		Short: "Hez makes orchestrating containers easy.",
		Long:  `Hez is a tool to orchestrate containers and manage the environment around it.`,
	}

	hezCli.AddCommand(apply.NewCommand())
	hezCli.AddCommand(proxy.NewCommand())
	hezCli.AddCommand(ssl.NewCommand())

	if err := hezCli.Execute(); err != nil {
		os.Exit(1)
	}
}
