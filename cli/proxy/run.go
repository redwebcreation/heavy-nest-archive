package proxy

import (
	"github.com/redwebcreation/hez/internal"
	"github.com/redwebcreation/hez/internal/proxy"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func RunRunCommand(_ *cobra.Command, _ []string) error {
	proxy := proxy.Proxy{
		Paths: proxy.PathsConfig{
			CertificatesDirectory: internal.CertificatesDirectory,
			KeyFile:               internal.DataDirectory + "/key.pem",
			CertFile:              internal.DataDirectory + "/cert.pem",
		},
		Http: proxy.HttpConfig{
			Port: internal.Config.Proxy.Http.Port,
		},
		Https: proxy.HttpsConfig{
			Port:       internal.Config.Proxy.Https.Port,
			SelfSigned: *internal.Config.Proxy.Https.SelfSigned,
		},
		Handler: func(writer http.ResponseWriter, request *http.Request) {
			application := internal.Config.Applications[request.Host]

			if application == nil {
				logWithRequestContext("no such host", request)
				writer.WriteHeader(503)
				_, _ = writer.Write([]byte("Service unavailable."))
				return
			}

			container, err := application.GetContainer(internal.AnyType)

			if err != nil {
				logWithRequestContext(
					"container missing",
					request,
					zap.String("container_name", application.Name(internal.ApplicationContainer)),
				)
				return
			}

			containerUrl, err := url.Parse("http://" + container.Ip + ":" + application.ContainerPort)

			if err != nil {
				logWithRequestContext("invalid url", request)
				return
			}

			// create the reverse proxy
			proxy := &httputil.ReverseProxy{
				Director: func(req *http.Request) {
					req.URL.Scheme = containerUrl.Scheme
					req.URL.Host = containerUrl.Host
					req.Header.Set("X-Forwarded-Host", application.Host)
					req.Header.Set("X-Forwarded-Proto", containerUrl.Scheme)
					req.Host = application.Host
				},
				ModifyResponse: func(response *http.Response) error {
					response.Header.Del("Server")

					return nil
				},
				ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
					writer.WriteHeader(http.StatusBadGateway)
					_, _ = writer.Write([]byte(http.StatusText(http.StatusBadGateway)))
					internal.Logger.Fatal(err.Error())
				},
			}

			// Note that ServeHttp is non-blocking and uses a go routine under the hood
			proxy.ServeHTTP(writer, request)

			logWithRequestContext("request served", request)
		},
		IsAllowedHost: func(host string) bool {
			for appHost := range internal.Config.Applications {
				if appHost == host {
					return true
				}
			}

			return false
		},
	}

	return proxy.Serve()
}

func RunCommand() *cobra.Command {
	command := internal.CreateCommand(&cobra.Command{
		Use:   "run",
		Short: "Starts the proxy server.",
		Long:  `Starts the proxy server on the configured ports.`,
	}, nil, RunRunCommand)

	return command
}

func logWithRequestContext(message string, request *http.Request, fields ...zap.Field) {
	ip := strings.Split(request.RemoteAddr, ":")[0]

	fields = append(fields, zap.String("vhost", request.Host))
	fields = append(fields, zap.String("method", request.Method))
	fields = append(fields, zap.String("request_uri", request.RequestURI))
	fields = append(fields, zap.String("ip", ip))
	fields = append(fields, zap.String("referer", request.Referer()))
	fields = append(fields, zap.String("ua", request.UserAgent()))

	internal.Logger.Info(
		message,
		fields...,
	)
}
