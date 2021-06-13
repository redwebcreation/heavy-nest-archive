package proxy

import (
	"errors"
	"github.com/redwebcreation/hez/ansi"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
)

func runEnableCommand(_ *cobra.Command, _ []string) error {
	if !core.IsRunningAsRoot() {
		return errors.New("this command requires elevated privileges")
	}

	if core.IsProxyEnabled() {
		return errors.New("proxy is already enabled")
	}

	err := core.EnableProxy()

	if err != nil {
		return err
	}

	ansi.Text("Proxy has been successfully enabled.", ansi.Green)
	return nil
}

func EnableCommand() *cobra.Command {
	return core.CreateCommand(&cobra.Command{
		Use:   "enable",
		Short: "Enables the reverse proxy.",
		Long:  `Registers the reverse proxy in systemd`,
	}, nil, runEnableCommand)
}
