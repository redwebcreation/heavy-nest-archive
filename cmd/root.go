package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hez",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Print(`
$$\                           
$$ |                          
$$$$$$$\   $$$$$$\  $$$$$$$$\ 
$$  __$$\ $$  __$$\ \____$$  |
$$ |  $$ |$$$$$$$$ |  $$$$ _/ 
$$ |  $$ |$$   ____| $$  _/   
$$ |  $$ |\$$$$$$$\ $$$$$$$$\ 
\__|  \__| \_______|\________|

`)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
