package proxy

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

func RunRunCommand(_ *cobra.Command, _ []string) error {
	http.HandleFunc("/", core.HandleRequest)

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			domains := core.GetWhitelistedDomains()

			for _, domain := range domains {
				if domain == host {
					return nil
				}
			}

			return errors.New("host not configured in whitelist")
		},
		Cache: autocert.DirCache(core.CertificatesDirectory),
	}

	go func() {
		// HTTP server that redirects to the HTTPS one.
		h := certManager.HTTPHandler(nil)
		err := http.ListenAndServe(":"+core.Config.Proxy.Http.Port, h)

		if err != nil {
			core.Logger.Fatal(err.Error())
		}
	}()

	server := &http.Server{
		Addr: ":" + core.Config.Proxy.Https.Port,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	err := server.ListenAndServeTLS("", "")

	if err != nil {
		core.Logger.Fatal(err.Error())
	}

	return nil
}

func RunCommand() *cobra.Command {
	command := core.CreateCommand(&cobra.Command{
		Use:   "run",
		Short: "Starts the proxy server.",
		Long:  `Starts the proxy server on the configured ports.`,
	}, nil, RunRunCommand)

	return command
}
