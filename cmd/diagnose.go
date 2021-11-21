package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/common"
	"github.com/wormable/ui"
)

func runDiagnoseCommand(_ *cobra.Command, _ []string) error {
	diagnosis := common.AnalyseConfig()

	fmt.Printf("\n  %sErrors:%s\n", ui.Red.Fg(), ui.Stop)
	if diagnosis.ErrorCount == 0 {
		fmt.Printf("  %s- no errors%s\n", ui.Gray.Fg(), ui.Stop)
	}
	for _, report := range diagnosis.Errors {
		fmt.Printf("  - %s\n", ui.White.Fg()+report.Title+ui.Stop)
		if report.Error != nil {
			fmt.Printf("  %s\n", ui.White.Fg()+ui.Dim+report.Error.Error()+ui.Stop)
		}
	}

	fmt.Printf("\n  %sWarnings:%s\n", ui.Yellow.Fg(), ui.Stop)
	if diagnosis.WarningCount == 0 {
		fmt.Printf("  %s- no warnings%s\n", ui.Gray.Fg(), ui.Stop)
	}
	for _, report := range diagnosis.Warnings {
		fmt.Printf("  - %s\n", ui.White.Fg()+report.Title+ui.Stop)
		fmt.Printf("  %s\n", ui.White.Fg()+ui.Dim+report.Advice+ui.Stop)
	}

	fmt.Printf(
		"\n  %sFound %s %d %s errors and %s%d%s warnings in the configuration.%s\n",
		ui.White.Fg(),
		If(diagnosis.ErrorCount > 0, ui.Red.Bg(), ui.Green.Bg()),
		diagnosis.ErrorCount,
		ui.Stop+ui.White.Fg(),
		If(diagnosis.WarningCount > 0, ui.Stop+ui.Yellow.Fg(), ui.Stop+ui.Green.Fg()),
		diagnosis.WarningCount,
		ui.Stop+ui.White.Fg(),
		ui.Stop,
	)

	return nil
}

func DiagnoseCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "diagnose",
		Short: "Display diagnostic information that helps you fix your config",
	}, runDiagnoseCommand, nil)
}

func If(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	} else {
		return b
	}
}
