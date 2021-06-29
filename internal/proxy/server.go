package proxy

import (
	"crypto/tls"
	"github.com/go-http-utils/etag"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HttpConfig struct {
	Port string
}

type HttpsConfig struct {
	Port       string
	SelfSigned bool
}

type PathsConfig struct {
	CertificatesDirectory string
	KeyFile               string
	CertFile              string
}

type Proxy struct {
	Paths         PathsConfig
	Http          HttpConfig
	Https         HttpsConfig
	Handler       func(writer http.ResponseWriter, request *http.Request)
	HostToUrl     func(host string) *url.URL
	IsAllowedHost func(host string) bool
}

func (p Proxy) RequestHandler() http.Handler {
	proxy := func(rw http.ResponseWriter, request *http.Request) {
		p.Handler(rw, request)
	}

	composer := NewComposer(proxy).
		Use(Compress).
		Use(func(handler http.Handler) http.Handler {
			return etag.Handler(handler, true)
		})
	return composer.Handler
}

func (p Proxy) Serve() error {
	certificateManager := p.getCertificateManager()

	http.Handle("/", p.RequestHandler())

	go func() {
		err := http.ListenAndServe(":"+p.Http.Port, certificateManager.HTTPHandler(http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, p.getSecureUrl(request), http.StatusPermanentRedirect)
		}))))

		if err != nil {
			log.Fatal(err)
		}
	}()

	if p.Https.SelfSigned {
		err := p.createSelfSignedCertificates()
		if err != nil {
			return err
		}

		err = http.ListenAndServeTLS(":"+p.Https.Port, p.Paths.CertFile, p.Paths.KeyFile, nil)
		return err
	}

	server := &http.Server{
		Addr: ":" + p.Https.Port,
		TLSConfig: &tls.Config{
			GetCertificate: certificateManager.GetCertificate,
		},
	}

	err := server.ListenAndServeTLS("", "")

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (p Proxy) getSecureUrl(request *http.Request) string {
	u := "https://" + request.Host
	u = strings.Replace(u, p.Http.Port, p.Https.Port, 1)

	if p.Https.Port == "443" {
		u = strings.Replace(u, ":"+p.Https.Port, "", 1)
	}

	u += request.URL.String()

	return u
}
