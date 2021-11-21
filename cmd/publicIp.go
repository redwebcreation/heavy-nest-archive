package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/ui"
	"io/ioutil"
	"net/http"
	"strings"
)

func runPublicIp(_ *cobra.Command, _ []string) error {
	response, err := http.Get("http://checkip.amazonaws.com")
	ui.Check(err)
	defer response.Body.Close()

	rawPublicIp, err := ioutil.ReadAll(response.Body)
	ui.Check(err)

	publicIp := strings.TrimSpace(string(rawPublicIp))
	fmt.Println(publicIp)
	return nil
}

func PublicIpCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "public-ip",
		Short: "Get this computer's public IP address",
	}, runPublicIp, nil)
}
