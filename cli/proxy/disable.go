package proxy

import (
	"errors"
	"fmt"
	"github.com/redwebcreation/hez/internal"
	"github.com/spf13/cobra"
)

func runDisableCommand(_ *cobra.Command, _ []string) error {
	err := internal.ElevateProcess()

	if err != nil {
		return err
	}

	if !internal.IsProxyEnabled() {
		return errors.New("proxy is already disabled")
	}

	err = internal.DisableProxy()

	if err != nil {
		return err
	}

	fmt.Println("Proxy has been successfully disabled.")

	return nil
}

func DisableCommand() *cobra.Command {
	return internal.CreateCommand(&cobra.Command{
		Use:   "disable",
		Short: "Disables the reverse proxy.",
		Long:  `Disables the reverse proxy configuration file in systemd`,
	}, nil, runDisableCommand)
}
