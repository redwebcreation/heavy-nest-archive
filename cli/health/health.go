package health

import (
	"github.com/spf13/cobra"
)

func run(_ *cobra.Command, _ []string) {
}

func NewCommand() *cobra.Command {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Returns the system's health",
		Run:   run,
	}

	return healthCmd
}
