package proxy

import (
	"github.com/redwebcreation/hez/core"
	"github.com/redwebcreation/hez/ntee"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func RunRunCommand(_ *cobra.Command, _ []string) error {
	proxy := ntee.Proxy{
		Paths: ntee.PathsConfig{
			CertificatesDirectory: core.CertificatesDirectory,
			KeyFile:               core.DataDirectory + "/key.pem",
			CertFile:              core.DataDirectory + "/cert.pem",
		},
		Http: ntee.HttpConfig{
			Port: core.Config.Proxy.Http.Port,
		},
		Https: ntee.HttpsConfig{
			Port:       core.Config.Proxy.Https.Port,
			SelfSigned: *core.Config.Proxy.Https.SelfSigned,
		},
		Handler: func(writer http.ResponseWriter, request *http.Request) {
			application := core.Config.Applications["madeinfranz.fr"]

			if application == nil {
				logWithRequestContext("no such host", request)
				writer.WriteHeader(503)
				_, _ = writer.Write([]byte("Service unavailable."))
				return
			}

			container, err := application.GetContainer(core.AnyType)

			if err != nil {
				logWithRequestContext(
					"container missing",
					request,
					zap.String("container_name", application.Name(core.ApplicationContainer)),
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
					core.Logger.Fatal(err.Error())
				},
			}

			// Note that ServeHttp is non blocking and uses a go routine under the hood
			proxy.ServeHTTP(writer, request)

			logWithRequestContext("request served", request)
		},
		IsAllowedHost: nil,
	}

	return proxy.Serve()
}

func RunCommand() *cobra.Command {
	command := core.CreateCommand(&cobra.Command{
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

	core.Logger.Info(
		message,
		fields...,
	)
}
