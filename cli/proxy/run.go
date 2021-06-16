package proxy

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func RunRunCommand(_ *cobra.Command, _ []string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		for _, application := range core.Config.Applications {
			if request.Host == application.Host {
				container, err := application.GetContainer(core.AnyType)

				if err != nil {
					core.Logger.Error(
						"container missing",
						zap.String("vhost", request.Host),
						zap.String("container_name", application.Name(core.ApplicationContainer)),
					)
					break
				}

				containerUrl, err := url.Parse("http://" + container.Ip + ":" + application.ContainerPort)

				if err != nil {
					core.Logger.Error(
						"invalid url",
						zap.String("error", err.Error()),
					)
					break
				}

				// create the reverse proxy
				proxy := httputil.NewSingleHostReverseProxy(containerUrl)

				// Update the headers to allow for SSL redirection
				request.URL.Host = containerUrl.Host
				request.URL.Scheme = containerUrl.Scheme
				request.Header.Set("X-Forwarded-Host", request.Header.Get("Host"))
				request.Host = containerUrl.Host

				// Note that ServeHttp is non blocking and uses a go routine under the hood
				proxy.ServeHTTP(writer, request)
			}
		}
	})

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
		//HTTP server that redirects to the HTTPS.
		err := http.ListenAndServe(":"+core.Config.Proxy.Http.Port, certManager.HTTPHandler(nil))

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
