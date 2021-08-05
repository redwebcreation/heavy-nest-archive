package proxy

import (
	"errors"
	"github.com/redwebcreation/hez/internal"
	ui "github.com/redwebcreation/hez/internal/ui"
	"github.com/spf13/cobra"
)

func runEnableCommand(_ *cobra.Command, _ []string) error {
	err := internal.ElevateProcess()

	if err != nil {
		return err
	}

	if internal.IsProxyEnabled() {
		return errors.New("proxy is already enabled")
	}

	err = internal.EnableProxy()

	if err != nil {
		return err
	}

	ui.Success("Proxy has been successfully enabled.")
	return nil
}

func EnableCommand() *cobra.Command {
	return internal.CreateCommand(&cobra.Command{
		Use:   "enable",
		Short: "Enables the reverse proxy.",
		Long:  `Registers the reverse proxy in systemd`,
	}, nil, runEnableCommand)
}
