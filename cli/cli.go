package cli

import (
	"fmt"

	"github.com/redwebcreation/hez/client"
)

func Execute() {
	fmt.Println()
	fmt.Println()
	deployment := client.DeploymentConfiguration{
		Name: "factures_redwebcreation_fr_80",
		Registry: &client.RegistryAuth{
			Username: "nologin",
			Password: "6fe186fd-edc2-485e-8f5d-93c357add27c",
		},
		Image: "rg.fr-par.scw.cloud/rwsapps/aaa:b",
		Host:  "facture.net",
		Port:  "80",
	}

	deployment.Deploy()
	//return
	/*
	   cli := &cobra.Command{
	   		Use:   "hez",
	   		Short: "Hez makes orchestrating containers easy.",
	   		Long:  "Hez is to tool to orchestrate containers and manage the environment around them.",
	   	}

	   	cli.SilenceErrors = true
	   	err := cli.Execute()new
	   	check(err)
	*/
}
