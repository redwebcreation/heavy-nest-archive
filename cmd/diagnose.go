package cmd

import (
	"fmt"
	"github.com/redwebcreation/nest/client"
	"github.com/redwebcreation/nest/internal"
	"github.com/redwebcreation/nest/ui"
	"github.com/spf13/cobra"
)

func runDiagnoseCommand(_ *cobra.Command, _ []string) error {
	diagnosis := client.Analyse(client.Config)

	if diagnosis.ErrorCount == 0 {
		fmt.Println()
		fmt.Printf("    %sRan %d checks.%s\n", ui.White.Fg(), len(diagnosis.Checks), ui.Stop)
		fmt.Printf("    %sNo errors found. Great job!%s\n", ui.Green.Fg(), ui.Stop)
		return nil
	}

	for _, diagnosis := range diagnosis.Errors {
		fmt.Printf("\n    - %s%s%s\n", ui.White.Fg(), diagnosis.Title, ui.Stop)
		if diagnosis.Error != nil {
			fmt.Printf("    %s%s%s\n ", ui.White.Fg()+ui.Dim, diagnosis.Error.Error(), ui.Stop)
		}
	}

	fmt.Printf("\n    %sFound %s %d %s errors in the configuration.%s\n", ui.White.Fg(), ui.Red.Bg(), diagnosis.ErrorCount, ui.Stop+ui.White.Fg(), ui.Stop)

	return nil
}

func DiagnoseCommand() *cobra.Command {
	return internal.CreateCommand(&cobra.Command{
		Use:   "diagnose",
		Short: "Display diagnostic information that helps you fix your config",
	}, nil, runDiagnoseCommand)
}
