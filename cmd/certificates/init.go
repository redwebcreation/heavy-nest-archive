package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
	"github.com/wormable/nest/globals"
	"github.com/wormable/ui"
	"math/big"
	"net"
	"os"
	"time"
)

var renew bool
var forceRenew bool

func runInitCommand(_ *cobra.Command, _ []string) error {
	_, err := os.Stat(globals.CACertificate)

	if !renew && !os.IsNotExist(err) {
		return fmt.Errorf("a root certificate already exists")
	}
	fmt.Printf("\n  Creating certificates in %s.\n\n", globals.CertificatesDirectory)

	if renew && !forceRenew {
		fmt.Printf("  | Renewing the certificates %swill invalidate all the backends connection%s\n", ui.White.Fg(), ui.Stop)
		fmt.Printf("  | This %swill induce downtime%s on this master and you will have to reconfigure all the backends.\n", ui.White.Fg(), ui.Stop)

		fmt.Println("  | You may use --force-renew to skip the waiting time and the message altogether.")
		fmt.Printf("  | Sleeping for 10 seconds, %splease abort if you have any doubts (Ctrl+C).%s\n\n", ui.Red.Fg(), ui.Stop)

		time.Sleep(10 * time.Second)
	}

	validity := 0 // in years
	cACertificate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Nest Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(validity, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	serverCertificate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(validity, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	ui.NewLog("Generating a new private key for the CA").Print()
	pkey := generateRsaKey(globals.CAPrivateKey)
	ui.NewLog("Generated a private key for the CA").Top(1).Print()

	ui.NewLog("Generating a new certificate for the CA").Print()
	generateCertificate(globals.CACertificate, cACertificate, cACertificate, &pkey.PublicKey, pkey)
	ui.NewLog("Generated a certificate for the CA").Top(1).Print()

	ui.NewLog("Generating a new private key for the Master").Print()
	certPrivateKey := generateRsaKey(globals.ServerPrivateKey)
	ui.NewLog("Generated a private key for the Master").Top(1).Print()

	ui.NewLog("Generating a new certificate for the Master").Print()
	generateCertificate(globals.ServerCertificate, serverCertificate, cACertificate, &certPrivateKey.PublicKey, pkey)
	ui.NewLog("Generated a certificate for the Master").Top(1).Print()
	return nil
}

func InitCommand() *cobra.Command {
	return cmd.Decorate(&cobra.Command{
		Use: "init",
	}, func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&renew, "renew", false, "Renew certificates")
		cmd.Flags().BoolVar(&forceRenew, "force-renew", false, "Force renew certificates")
	}, runInitCommand)
}

func generateRsaKey(name string) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	ui.Check(err)

	save(name, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return key
}

func generateCertificate(name string, template, parent *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
	ui.Check(err)

	save(name, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
}

func save(name string, block *pem.Block) {
	pkeyFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0600)
	ui.Check(err)

	err = pem.Encode(pkeyFile, block)
	ui.Check(err)
}
