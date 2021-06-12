package proxy

import (
	"crypto/tls"
	"github.com/redwebcreation/hez2/core"
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strconv"
)

func RunRunCommand(_ *cobra.Command, _ []string) error {
	lastApplyExecution := core.GetLastApplyExecution()
	proxiables, err := core.GetProxiableContainers()

	if err != nil {
		return err
	}

	http.HandleFunc("/", core.HandleRequest(lastApplyExecution, proxiables))

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(core.GetWhitelistedDomains()...),
		Cache:      autocert.DirCache(globals.CertificatesDirectory),
	}

	if globals.Config.Proxy.Http.Enabled {
		go func() {
			// HTTP server that redirects to the HTTPS one.
			h := certManager.HTTPHandler(nil)
			err := http.ListenAndServe(":"+strconv.Itoa(globals.Config.Proxy.Http.Port), h)

			if err != nil {
				globals.Logger.Fatal(err.Error())
			}
		}()
	}

	if globals.Config.Proxy.Https.Enabled {
		server := &http.Server{
			Addr: ":" + strconv.Itoa(globals.Config.Proxy.Https.Port),
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		err = server.ListenAndServeTLS("", "")

		if err != nil {
			globals.Logger.Fatal(err.Error())
		}
	}

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
