package proxy

import (
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
)

func RunRunCommand(_ *cobra.Command, _ []string) error {
	globals.Ansi.ErrorBlock("well there is still some work to do")
	// https://rolandshoemaker.github.io/acme-tls-alpn/draft-ietf-acme-tls-alpn.html
	// https://github.com/letsencrypt/pebble/
	// https://github.com/mdebski/golang-alpn-example
	return nil
}

func RunCommand() *cobra.Command {
	command := util.CreateCommand(&cobra.Command{
		Use:   "run",
		Short: "Starts the proxy server.",
		Long:  `Starts the proxy server on the configured ports.`,
	}, nil, RunRunCommand)

	return command
}
