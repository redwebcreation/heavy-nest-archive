package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"net"
	"net/http"
	"os/exec"
	"strconv"
)

var selfSigned bool
var Port int
var Ssl int

func runRunCommand(_ *cobra.Command, _ []string) {
	proxiableContainers, err := core.GetProxiableContainers()

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(proxiableContainers) == 0 {
		fmt.Println("Found 0 proxiable containers. Aborting.")
		return
	}

	handler := func(writer http.ResponseWriter, request *http.Request) {
		ip, _, _ := net.SplitHostPort(request.RemoteAddr)
		request.Header.Set("X-Forwarded-For", ip)

		for _, proxiableContainer := range proxiableContainers {
			//if request.Host == proxiableContainer.VirtualHost {
				success := core.ForwardRequest(proxiableContainer, writer, request)
				if success {
					zap.L().Info(
						"request.success",
						zap.String("method", request.Method),
						zap.String("ip", ip),
						zap.String("vhost", request.Host),
					)
				}
				return
			//}
		}

		zap.L().Info(
			"request.invalid",
			zap.String("method", request.Method),
			zap.String("ip", ip),
			zap.String("vhost", request.Host),
		)
		writer.WriteHeader(404)
		writer.Write([]byte("404. That’s an error. \nThe requested URL " + request.RequestURI + " was not found on this server. That’s all we know."))
	}
	http.HandleFunc("/", handler)

	if selfSigned {
		go func() {
			err := http.ListenAndServe(":"+strconv.Itoa(Port), nil)

			if err != nil {
				zap.L().Fatal(err.Error())
			}
		}()

		keyPath, certPath := handleSSLForTesting()

		err := http.ListenAndServeTLS(":"+strconv.Itoa(Ssl), certPath, keyPath, nil)

		if err != nil {
			zap.L().Fatal(err.Error())
			return
		}
	} else {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(core.GetWhitelistedDomains()...),
		}

		server := &http.Server{
			Addr: ":" + strconv.Itoa(Ssl),
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go func() {
			// HTTP server that redirects to the HTTPS one.
			h := certManager.HTTPHandler(nil)
			err := http.ListenAndServe(":"+strconv.Itoa(Port), h)

			if err != nil {
				zap.L().Fatal(err.Error())
			}
		}()

		err := server.ListenAndServeTLS("", "")

		if err != nil {
			zap.L().Fatal(err.Error())
		}
	}
}

func initRunCommand() *cobra.Command {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Starts the proxy server",
		Run:   runRunCommand,
	}

	config, _ := core.FindConfig(core.ConfigFile()).Resolve()

	runCommand.Flags().BoolVar(&selfSigned, "self-signed", *config.Proxy.SelfSigned, "Generate self signed SSL certificates.")
	runCommand.Flags().IntVar(&Port, "port", config.Proxy.Port, "Runs the HTTP proxy on a specific port.")
	runCommand.Flags().IntVar(&Ssl, "ssl", config.Proxy.Ssl, "Runs the HTTPS proxy on a specific port. ")

	return runCommand
}

func handleSSLForTesting() (string, string) {
	keyPath := core.StorageDirectory() + "/key.pem"
	certPath := core.StorageDirectory() + "/cert.pem"

	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyPath, "-out", certPath, "-days", "365", "-nodes", "-subj", "/CN=localhost")

	_ = cmd.Run()

	return keyPath, certPath
}
