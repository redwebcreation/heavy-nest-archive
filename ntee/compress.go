package ntee

import (
	"compress/flate"
	"compress/gzip"
	"github.com/go-http-utils/negotiator"
	"io"
	"mime"
	"net/http"
	"regexp"

	"github.com/GitbookIO/mimedb"
)

var compressibleTypeRegExp = regexp.MustCompile(`(?i)^text/|\+json$|\+text$|\+xml$`)

func IsCompressible(contentType string) bool {
	dbMatched := false

	extensions, err := mime.ExtensionsByType(contentType)

	if err != nil {
		return false
	}

	for _, ext := range extensions {
		// extensions all start with a dot
		if entry, ok := mimedb.DB[ext[1:]]; ok {
			dbMatched = true

			if entry.Compressible {
				return true
			}
		}
	}

	if !dbMatched && compressibleTypeRegExp.MatchString(contentType) {
		return true
	}

	return false
}

type compressWriter struct {
	rw http.ResponseWriter
	w  io.WriteCloser
}

func (cw compressWriter) Header() http.Header {
	return cw.rw.Header()
}

func (cw compressWriter) WriteHeader(status int) {
	cw.rw.Header().Del(ContentLength)

	cw.rw.WriteHeader(status)
}

func (cw compressWriter) Write(b []byte) (int, error) {
	cw.rw.Header().Del(ContentLength)

	return cw.w.Write(b)
}

// Compress wraps the http.Handler h with compress support.
func Compress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var w io.WriteCloser

		if req.Method != http.MethodHead &&
			res.Header().Get(ContentEncoding) == "" &&
			IsCompressible(req.Header.Get(ContentType)) {

			n := negotiator.New(req.Header)
			encoding := n.Encoding("gzip", "deflate")

			switch encoding {
			case "gzip":
				w = gzip.NewWriter(res)
			case "deflate":
				w, _ = flate.NewWriter(res, flate.DefaultCompression)
			}

			cw := compressWriter{rw: res, w: w}

			cw.Header().Set(ContentEncoding, encoding)
			cw.Header().Set(Vary, AcceptEncoding)

			defer cw.w.Close()

			h.ServeHTTP(cw, req)
			return
		}

		h.ServeHTTP(res, req)
	})
}
