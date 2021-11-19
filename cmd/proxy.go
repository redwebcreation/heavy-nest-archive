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
	"regexp"
)

func runProxyCommand(_ *cobra.Command, _ []string) error {
	certificateManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if common.Config.Applications[host].Host != "" {
				return nil
			}

			// TODO: Here actualize the config and check again
			return fmt.Errorf("host is not allowed")
		},
		Cache: autocert.DirCache(common.Config.Production.CertificateCache),
	}

	go func() {
		err := http.ListenAndServe(":"+common.Config.Production.HttpPort, certificateManager.HTTPHandler(nil))

		common.Config.Log(syslog.LOG_ERR, err.Error())
	}()

	if common.Config.Production.SelfSigned {
		panic("not implemented")
	}

	server := &http.Server{
		Addr: common.Config.Production.HttpsPort,
		TLSConfig: &tls.Config{
			GetCertificate: certificateManager.GetCertificate,

		},
	}

	return server.ListenAndServeTLS("", "")
	//server := http.Server{
	//	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		host := r.Host
	//
	//		application := common.Config.Applications[host]
	//
	//		// application is not its null version
	//		if application.Host != host {
	//			w.WriteHeader(http.StatusNotFound)
	//			return
	//		}
	//
	//		url, err := application.Url()
	//
	//		if err != nil {
	//			if err.Error() == "dead proxy" {
	//				w.WriteHeader(http.StatusBadGateway)
	//				_, _ = w.Write([]byte("Proxy died."))
	//			} else {
	//				w.WriteHeader(http.StatusInternalServerError)
	//			}
	//			return
	//		}
	//
	//		httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
	//	}),
	//}
	//
	//return server.ListenAndServe()
	return nil
}
func ProxyCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use: "proxy",
	}, nil, runProxyCommand)
}
func getSecureUrl(r *http.Request) string {
	return "https://" + (regexp.MustCompile("/(:[0-9]+)/").ReplaceAllString(r.Host, "")) + r.URL.Path
}

func logRequest(r *http.Request, priority syslog.Priority, message string) {
	context := common.Fields{
		"ip":     getRequestIp(r),
		"ua":     r.UserAgent(),
		"scheme": r.URL.Scheme,
		"host":   r.Host,
		"path":   r.URL.Path,
	}
	for _, policy := range common.Config.LogPolicies {
		policy.WithContext(context).Log(priority, message)
	}
}

func getRequestIp(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
