package proxy

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/klauspost/compress/gzhttp"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
)

func GetCertificateManager() autocert.Manager {
	return autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			for _, application := range core.Config.Applications {
				if host == application.Host {
					return nil
				}
			}

			return errors.New("host not configured in whitelist")
		},
		Cache: autocert.DirCache(core.CertificatesDirectory),
	}

}

func RequestHandler(writer http.ResponseWriter, request *http.Request) {
	// TODO: Use Host as the key to remove this loop
	var requestServed bool

	for _, application := range core.Config.Applications {
		if request.Host == application.Host {
			requestServed = true
			container, err := application.GetContainer(core.AnyType)

			if err != nil {
				logWithRequestContext(
					"container missing",
					request,
					zap.String("container_name", application.Name(core.ApplicationContainer)),
				)
				break
			}

			containerUrl, err := url.Parse("http://" + container.Ip + ":" + application.ContainerPort)

			if err != nil {
				logWithRequestContext("invalid url", request)
				break
			}

			// create the reverse proxy
			proxy := httputil.NewSingleHostReverseProxy(containerUrl)

			// Update the headers to allow for SSL redirection
			request.URL.Host = containerUrl.Host
			request.URL.Scheme = containerUrl.Scheme
			request.Header.Set("X-Forwarded-Host", request.Header.Get("Host"))
			request.Header.Set("X-Forwarded-Proto", request.URL.Scheme)
			request.Host = application.Host

			// Note that ServeHttp is non blocking and uses a go routine under the hood
			proxy.ServeHTTP(writer, request)

			logWithRequestContext("request served", request)
		}
	}

	if !requestServed {
		logWithRequestContext("request washed up", request)
		writer.WriteHeader(503)
		_, _ = writer.Write([]byte("Service unavailable."))
	}
}

func RunRunCommand(_ *cobra.Command, _ []string) error {
	certificateManager := GetCertificateManager()
	handler := gzhttp.GzipHandler(http.HandlerFunc(RequestHandler))
	http.Handle("/", handler)

	go func() {
		err := http.ListenAndServe(
			":"+core.Config.Proxy.Http.Port,
			certificateManager.HTTPHandler(nil),
		)

		if err != nil {
			core.Logger.Fatal(err.Error())
		}
	}()

	if *core.Config.Proxy.Https.SelfSigned {
		keyPath, certPath, err := createCertificates()

		if err != nil {
			return err
		}

		err = http.ListenAndServeTLS(":"+core.Config.Proxy.Https.Port, certPath, keyPath, nil)

		if err != nil {
			core.Logger.Fatal(err.Error())
		}
	} else {
		server := &http.Server{
			Addr: ":" + core.Config.Proxy.Https.Port,
			TLSConfig: &tls.Config{
				GetCertificate: certificateManager.GetCertificate,
			},
		}
		err := server.ListenAndServeTLS("", "")

		if err != nil {
			core.Logger.Fatal(err.Error())
		}
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
func createCertificates() (string, string, error) {
	keyPath := core.DataDirectory + "/key.pem"
	certPath := core.DataDirectory + "/cert.pem"

	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyPath, "-out", certPath, "-days", "365", "-nodes", "-subj", "/CN=localhost")

	err := cmd.Run()

	if err != nil {
		return "", "", err
	}

	return keyPath, certPath, nil
}

func logWithRequestContext(message string, request *http.Request, fields ...zap.Field) {
	fields = append(fields, zap.String("vhost", request.Host))
	fields = append(fields, zap.String("ua", request.UserAgent()))
	fields = append(fields, zap.String("referer", request.Referer()))
	fields = append(fields, zap.String("method", request.Method))
	fields = append(fields, zap.String("request_uri", request.RequestURI))
	fields = append(fields, zap.String("ip", request.RemoteAddr))

	core.Logger.Info(
		message,
		fields...,
	)
}
