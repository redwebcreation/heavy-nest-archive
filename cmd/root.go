package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hez",
	Short: "Hez makes orchestrating containers easy.",
	Long: `Hez is a tool to orchestrate containers and manage the environment around it.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
