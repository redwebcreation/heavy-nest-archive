package ssl

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	sslCommand := &cobra.Command{
		Use: "ssl",
	}

	sslCommand.AddCommand(initGenerateCommand())

	return sslCommand
}
