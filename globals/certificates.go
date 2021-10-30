package globals

import (
	"github.com/mitchellh/go-homedir"
	"github.com/wormable/ui"
	"os"
)

var CertificatesDirectory string
var CACertificate string
var CAPrivateKey string
var ServerCertificate string
var ServerPrivateKey string

func init() {
	home, err := homedir.Dir()
	ui.Check(err)

	CertificatesDirectory = home + "/certs/"
	err = os.MkdirAll(CertificatesDirectory, os.FileMode(0777))
	ui.Check(err)

	CACertificate = CertificatesDirectory + "ca-cert.pem"
	CAPrivateKey = CertificatesDirectory  + "ca-key.pem"

	ServerCertificate = CertificatesDirectory + "server-cert.pem"
	ServerPrivateKey = CertificatesDirectory + "server-key.pem"
}
