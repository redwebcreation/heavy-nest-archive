package main

import (
	"fmt"
	"github.com/redwebcreation/hez/cli/apply"
	"github.com/redwebcreation/hez/cli/proxy"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	fmt.Print(`$$\                           
$$ |                          
$$$$$$$\   $$$$$$\  $$$$$$$$\ 
$$  __$$\ $$  __$$\ \____$$  |
$$ |  $$ |$$$$$$$$ |  $$$$ _/ 
$$ |  $$ |$$   ____| $$  _/   
$$ |  $$ |\$$$$$$$\ $$$$$$$$\
\__|  \__| \_______|\________|

`)
	hezCli := &cobra.Command{
		Use:   "hez",
		Short: "Hez makes orchestrating containers easy.",
		Long:  `Hez is a tool to orchestrate containers and manage the environment around it.`,
	}

	hezCli.AddCommand(apply.NewCommand())
	hezCli.AddCommand(proxy.NewCommand())

	if err := hezCli.Execute(); err != nil {
		os.Exit(1)
	}
}
