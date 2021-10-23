package globals

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

var certifaceAuthority string

func init() {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Organization: []string["Nest"],
			Country: []string{""},
			Province: []string{""},
			Locality: []string{""},
			StreetAddress: []string{""},
			PostalCode: []string{""},
		},
		NotBefore: time.Now(),
		NotAfter: time.Now().AddDate(100),
		IsCA: true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}
