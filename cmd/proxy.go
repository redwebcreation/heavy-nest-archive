package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wormable/nest/common"
	"net/http"
	"net/http/httputil"
)

func runProxyCommand(_ *cobra.Command, _ []string) error {
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			if host == common.Config.Staging.Host {
				w.WriteHeader(http.StatusNotImplemented)
				_, _ = w.Write([]byte("Not implemented"))
				return
			}

			application := common.Config.Applications[host]

			// application is not its null version
			if application.Host != host {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			url, err := application.Url()

			if err != nil {
				if err.Error() == "dead proxy" {
					w.WriteHeader(http.StatusBadGateway)
					_, _ = w.Write([]byte("Proxy died."))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
		}),
	}

	return server.ListenAndServe()
}
func ProxyCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use: "proxy",
	}, nil, runProxyCommand)
}
