package ssl

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var selfSigned bool

func runGenerateCommand(_ *cobra.Command, _ []string) {
	if !selfSigned {
		fmt.Println("Let's Encrypt is not supported yet.")
		os.Exit(1)
	}

	sslPath := core.ConfigDirectory() + "/ssl"

	_, err := os.Stat(sslPath)

	if os.IsNotExist(err) {
		os.Mkdir(sslPath, os.FileMode(0755))
	}

	keyPath := sslPath + "/key.pem"
	certPath := sslPath + "/cert.pem"

	os.Remove(keyPath)
	os.Remove(certPath)

	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyPath, "-out", certPath, "-days", "365", "-nodes", "-subj", "/CN=localhost")

	err = cmd.Run()

	fmt.Println("Creating a certificate at [" + certPath + "].")
	fmt.Println("Creating a key at [" + keyPath + "].")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func initGenerateCommand() *cobra.Command {
	generateCommand := &cobra.Command{
		Use:   "generate",
		Short: "Generates fresh SSL certificates.",
		Run:   runGenerateCommand,
	}

	generateCommand.Flags().BoolVar(&selfSigned, "self-signed", false, "Should the SSL certificate be self signed")

	return generateCommand
}
