package proxy

import (
	"errors"
	"fmt"
	"github.com/redwebcreation/hez2/core"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
)

func runDisableCommand(_ *cobra.Command, _ []string) error {
	if !core.IsRunningAsRoot() {
		return errors.New("this command requires elevated privileges")
	}

	if !core.IsProxyEnabled() {
		return errors.New("proxy is already disabled")
	}

	err := core.DisableProxy()

	if err != nil {
		return err
	}

	fmt.Println("Proxy has been successfully disabled.")

	return nil
}

func DisableCommand() *cobra.Command {
	return util.CreateCommand(&cobra.Command{
		Use:   "disable",
		Short: "Disables the reverse proxy.",
		Long:  `Disables the reverse proxy configuration file in systemd`,
	}, nil, runDisableCommand)
}
