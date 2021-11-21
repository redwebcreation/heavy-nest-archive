package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/common"
	"github.com/wormable/nest/ansi")

func runDiagnoseCommand(_ *cobra.Command, _ []string) error {
	diagnosis := common.AnalyseConfig()

	fmt.Printf("\n  %sErrors:%s\n", ansi.Red.Fg(), ansi.Reset)
	if diagnosis.ErrorCount == 0 {
		fmt.Printf("  %s- no errors%s\n", ansi.Gray.Fg(), ansi.Reset)
	}
	for _, report := range diagnosis.Errors {
		fmt.Printf("  - %s\n", ansi.White.Fg()+report.Title+ansi.Reset)
		if report.Error != nil {
			fmt.Printf("  %s\n", ansi.White.Fg()+ansi.Dim+report.Error.Error()+ansi.Reset)
		}
	}

	fmt.Printf("\n  %sWarnings:%s\n", ansi.Yellow.Fg(), ansi.Reset)
	if diagnosis.WarningCount == 0 {
		fmt.Printf("  %s- no warnings%s\n", ansi.Gray.Fg(), ansi.Reset)
	}
	for _, report := range diagnosis.Warnings {
		fmt.Printf("  - %s\n", ansi.White.Fg()+report.Title+ansi.Reset)
		fmt.Printf("  %s\n", ansi.White.Fg()+ansi.Dim+report.Advice+ansi.Reset)
	}

	fmt.Printf(
		"\n  %sFound %s %d %s errors and %s%d%s warnings in the configuration.%s\n",
		ansi.White.Fg(),
		If(diagnosis.ErrorCount > 0, ansi.Red.Bg(), ansi.Green.Bg()),
		diagnosis.ErrorCount,
		ansi.Reset+ansi.White.Fg(),
		If(diagnosis.WarningCount > 0, ansi.Reset+ansi.Yellow.Fg(), ansi.Reset+ansi.Green.Fg()),
		diagnosis.WarningCount,
		ansi.Reset+ansi.White.Fg(),
		ansi.Reset,
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
