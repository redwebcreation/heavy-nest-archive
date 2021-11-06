package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/wormable/nest/globals"
	"github.com/wormable/ui"
	"math/big"
	"net"
	"os"
	"time"
)

func CA() (*x509.Certificate, *rsa.PrivateKey) {
	bytes, err := os.ReadFile(globals.ServerCertificate)
	ui.Check(err)
	block, _ := pem.Decode(bytes)
	certificate, err := x509.ParseCertificate(block.Bytes)
	ui.Check(err)

	bytes, err = os.ReadFile(globals.ServerPrivateKey)
	ui.Check(err)
	key, err := x509.ParsePKCS1PrivateKey(bytes)

	return certificate, key
}

func CreateCertificate(host net.IP) {
	certificate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		IPAddresses:  []net.IP{host},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

}

func GenerateRsaKey(name string) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	ui.Check(err)

	save(name, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return key
}

func GenerateCertificate(name string, template, parent *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
	ui.Check(err)

	save(name, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
}

func save(name string, block *pem.Block) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0600)
	ui.Check(err)

	err = pem.Encode(file, block)
	ui.Check(err)
}
