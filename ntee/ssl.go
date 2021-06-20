package ntee

import (
	"bytes"
	"context"
	"errors"
	"golang.org/x/crypto/acme/autocert"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (p Proxy) getCertificateManager() autocert.Manager {
	return autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if p.IsAllowedHost(host) {
				return nil
			}

			return errors.New("host is not allowed")
		},
		Cache: autocert.DirCache(p.Paths.CertificatesDirectory),
	}

}
func (p Proxy) createSelfSignedCertificates() error {
	keyFile := p.Paths.KeyFile
	certFile := p.Paths.CertFile

	_, keyExists := os.Stat(keyFile)
	_, certExists := os.Stat(certFile)

	if !os.IsNotExist(keyExists) && !os.IsNotExist(certExists) {
		return nil
	}

	keyPath, _ := filepath.Abs(strings.TrimRight(keyFile, "/") + "/..")
	certPath, _ := filepath.Abs(strings.TrimRight(certFile, "/") + "/..")

	_ = os.MkdirAll(keyPath, os.FileMode(0777))
	_ = os.MkdirAll(certPath, os.FileMode(0777))

	var stderr bytes.Buffer

	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyFile, "-out", certFile, "-days", "365", "-nodes", "-subj", "/CN=localhost")

	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil && stderr.Len() > 0 {
		return errors.New(stderr.String())
	}

	return err
}
