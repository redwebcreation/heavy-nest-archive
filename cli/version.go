package cli

import (
	"fmt"
	"github.com/me/nest/global"
	"github.com/spf13/cobra"
)

func VersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Printf("nest@%s\n", global.Version)
			return nil
		},
	}
}