package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/common"
	"golang.org/x/crypto/acme/autocert"
	"log/syslog"
	"net/http"
	"net/http/httputil"
)

func runProxyCommand(_ *cobra.Command, _ []string) error {
	certificateManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if common.Config.Applications[host].Host != "" {
				return nil
			}

			// retry in case the config has been updated
			common.LoadConfig()

			if common.Config.Applications[host].Host != "" {
				return nil
			}

			return fmt.Errorf("host is not allowed")
		},
		Cache: autocert.DirCache(common.Config.Proxy.CertificateCache),
	}

	go func() {
		err := http.ListenAndServe(":"+common.Config.Proxy.HttpPort, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			certificateManager.HTTPHandler(nil).ServeHTTP(writer, request)
			logRequest(request, syslog.LOG_INFO, "to https")

		}))

		common.Config.Log(syslog.LOG_ERR, err.Error())
	}()

	if common.Config.Proxy.SelfSigned {
		panic("not implemented")
	}

	server := &http.Server{
		Addr: common.Config.Proxy.HttpsPort,
		TLSConfig: &tls.Config{
			GetCertificate: certificateManager.GetCertificate,
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			application := common.Config.Applications[host]

			// application exists
			if application.Host != host {
				w.WriteHeader(http.StatusNotFound)
				logRequest(r, syslog.LOG_INFO, "not found")
				return
			}

			url, err := application.Url()

			if err != nil {
				if err.Error() == "dead proxy" {
					w.WriteHeader(http.StatusBadGateway)
					_, _ = w.Write([]byte("Proxy died."))
					logRequest(r, syslog.LOG_ERR, "bad gateway")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					logRequest(r, syslog.LOG_ERR, err.Error())
				}
				return
			}

			httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)

			logRequest(r, syslog.LOG_INFO, "served")
		}),
	}

	return server.ListenAndServeTLS("", "")
}
func ProxyCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "proxy",
		Short: "start a proxy forwarding to your applications",
	}, runProxyCommand, nil)
}

func logRequest(r *http.Request, priority syslog.Priority, message string) {
	c := common.Fields{
		"ip":     getRequestIp(r),
		"ua":     r.UserAgent(),
		"scheme": r.URL.Scheme,
		"host":   r.Host,
		"path":   r.URL.Path,
	}
	for _, policy := range common.Config.LogPolicies {
		policy.WithContext(c).Log(priority, message)
	}
}

func getRequestIp(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
