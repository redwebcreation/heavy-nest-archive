package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

var selfSigned bool
var Port string
var Ssl string

func runRunCommand(_ *cobra.Command, _ []string) {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		proxiableContainers := core.GetProxiableContainers()
		logger := core.Logger()

		if len(proxiableContainers) == 0 {
			fmt.Println("Found 0 proxiable containers. Aborting.")
			os.Exit(1)
		}

		ip, _, _ := net.SplitHostPort(request.RemoteAddr)
		request.Header.Set("X-Forwarded-For", ip)

		for _, proxiableContainer := range proxiableContainers {
			if request.Host == proxiableContainer.VirtualHost {
				logger.Info(
					"request.handled",
					zap.String("method", request.Method),
					zap.String("ip", ip),
					zap.String("vhost", proxiableContainer.VirtualHost),
					zap.String("request_uri", request.RequestURI),
				)

				core.ForwardRequest(proxiableContainer, writer, request)
				return
			}
		}

		logger.Info(
			"request.invalid",
			zap.String("method", request.Method),
			zap.String("ip", ip),
			zap.String("vhost", request.Host),
			zap.String("request_uri", request.RequestURI),
		)
		writer.WriteHeader(404)
		writer.Write([]byte("404. That’s an error. \nThe requested URL " + request.RequestURI + " was not found on this server. That’s all we know."))
	}
	http.HandleFunc("/", handler)

	if selfSigned {
		go func() {
			log.Fatal(http.ListenAndServe(":"+Port, nil))
		}()

		keyPath, certPath := handleSSLForTesting()

		err := http.ListenAndServeTLS(":"+Ssl, certPath, keyPath, nil)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(core.GetWhitelistedDomains()...),
		}

		server := &http.Server{
			Addr: ":" + Ssl,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go func() {
			// HTTP server that redirects to the HTTPS one.
			h := certManager.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":"+Port, h))
		}()

		log.Fatal(server.ListenAndServeTLS("", ""))
	}
}

func initRunCommand() *cobra.Command {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Starts the proxy server",
		Run:   runRunCommand,
	}

	config := core.GetConfig()

	runCommand.Flags().BoolVar(&selfSigned, "self-signed", config.Proxy.SelfSigned, "Generate self signed SSL certificates.")
	runCommand.Flags().StringVar(&Port, "port", config.Proxy.Port, "Runs the HTTP proxy on a specific port.")
	runCommand.Flags().StringVar(&Ssl, "ssl", config.Proxy.Ssl, "Runs the HTTPS proxy on a specific port. ")

	return runCommand
}

func handleSSLForTesting() (string, string) {
	keyPath := core.StorageDirectory() + "/key.pem"
	certPath := core.StorageDirectory() + "/cert.pem"

	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyPath, "-out", certPath, "-days", "365", "-nodes", "-subj", "/CN=localhost")

	_ = cmd.Run()

	return keyPath, certPath
}
