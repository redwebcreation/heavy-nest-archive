package proxy

import (
	"bytes"
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func runStatusCommand(_ *cobra.Command, _ []string) {
	isProxyEnabled, err := core.IsProxyEnabled()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !isProxyEnabled {
		fmt.Println("Proxy has been disabled.")
		os.Exit(1)
	}

	var stdOut bytes.Buffer
	cmd := exec.Command("systemctl", "status", "hezproxy")
	cmd.Stdout = &stdOut

	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf(stdOut.String())
}

func initStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Returns the status of the reverse proxy.",
		Run:   runStatusCommand,
	}
}
