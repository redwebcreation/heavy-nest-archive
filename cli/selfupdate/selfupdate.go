package selfupdate

import (
	"github.com/spf13/cobra"
)

func run(_ *cobra.Command, _ []string) {

}

func NewCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "self-update",
		Short: "Updates Hez to the latest version.",
		Run:   run,
	}

	return applyCmd
}
