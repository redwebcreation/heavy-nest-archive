package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/redwebcreation/hez2/core"
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

func RunRunCommand(_ *cobra.Command, _ []string) error {
	lastApplyExecution := core.GetLastApplyExecution()

	proxiables, err := core.GetProxiableContainers()

	if err != nil {
		return err
	}

	http.HandleFunc("/", core.HandleRequest(lastApplyExecution, proxiables))

	fmt.Println("HTTPS init.")
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(core.GetWhitelistedDomains()...),
	}

	fmt.Println("Cert manager created with domains : ")
	fmt.Println(core.GetWhitelistedDomains())
	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	fmt.Println("Server created.")

	go func() {
		fmt.Println("HTTP server (in the goroutine)")
		// HTTP server that redirects to the HTTPS one.
		h := certManager.HTTPHandler(nil)
		err := http.ListenAndServe(":80", h)

		if err != nil {
			globals.Logger.Fatal(err.Error())
		}
	}()
	fmt.Println("Serving TlS")
	err = server.ListenAndServeTLS("", "")
	fmt.Println("Server started")
	if err != nil {
		globals.Logger.Fatal(err.Error())
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
