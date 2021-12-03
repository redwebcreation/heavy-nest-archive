package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd/api/tokens"
	"github.com/wormable/nest/common"
	"golang.org/x/crypto/acme/autocert"
	"log/syslog"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
)

func runProxyCommand(_ *cobra.Command, _ []string) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host

		application := common.Config.Applications[host]

		// application exists
		if application.Host != host {
			if common.Config.ApiHost == host {
				handleApiRequest(w, r)
				logRequest(r, syslog.LOG_INFO, "served api")
				return
			}
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
	})
	certificateManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if common.Config.Applications[host].Host != "" || common.Config.ApiHost == host {
				return nil
			}

			// retry in case the config has been updated
			common.LoadConfig()

			if common.Config.Applications[host].Host != "" {
				return nil
			}

			return fmt.Errorf("host is not allowed")
		},
		Cache: autocert.DirCache(common.CertificateDirectory),
	}

	go func() {
		err := http.ListenAndServe(":"+common.Config.Proxy.HttpPort, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			certificateManager.HTTPHandler(nil).ServeHTTP(writer, request)
			logRequest(request, syslog.LOG_INFO, "to https")

		}))

		common.Config.Log(syslog.LOG_ERR, err.Error())
	}()

	if common.Config.Proxy.SelfSigned {
		keyFile, certFile, err := createSelfSignedCertificates()
		if err != nil {
			return err
		}

		http.ListenAndServeTLS(":"+common.Config.Proxy.HttpsPort, certFile, keyFile, handler)
	}

	server := &http.Server{
		Addr: ":" + common.Config.Proxy.HttpsPort,
		TLSConfig: &tls.Config{
			GetCertificate: certificateManager.GetCertificate,
		},
		Handler: handler,
	}

	return server.ListenAndServeTLS("", "")
}

func handleApiRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/deploy" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not supported"))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Unsupported media type"))
		return
	}

	token := r.Header.Get("Authorization")

	if !tokens.Token(token).Exists() {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	r.ParseForm()

	host := r.Form.Get("host")

	if host == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Host is missing"))
		return
	}

	application := common.Config.Applications[host]

	if application.Host != host {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Application not found"))
		return
	}

	go func() {
		application.Deploy(common.DeploymentOptions{
			Pull:         true,
			Healthchecks: !skipHealthchecks,
			Name:         application.TemporaryContainerName(),
		})
		application.Deploy(common.DeploymentOptions{
			Pull:         false,
			Healthchecks: !skipHealthchecks,
			Name:         application.ContainerName(),
		})

		application.StopContainer(application.TemporaryContainerName())
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deployment started."))
	return
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

func createSelfSignedCertificates() (string, string, error) {
	keyFile := filepath.Join(common.CertificateDirectory, "key.pem")
	certFile := filepath.Join(common.CertificateDirectory, "cert.pem")

	if _, err := os.Stat(keyFile); err == nil {
		return "", "", nil
	}

	if _, err := os.Stat(certFile); err == nil {
		return "", "", nil
	}

	var stderr bytes.Buffer

	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyFile, "-out", certFile, "-days", "365", "-nodes", "-subj", "/CN=localhost")

	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil && stderr.Len() > 0 {
		return "", "", errors.New(stderr.String())
	}

	return keyFile, certFile, err
}
