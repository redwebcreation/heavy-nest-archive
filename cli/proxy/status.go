package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os/exec"
)

func runStatusCommand(_ *cobra.Command, _ []string) error {
	if !core.IsProxyEnabled() {
		return errors.New("proxy is disabled")
	}

	var stdOut bytes.Buffer
	cmd := exec.Command("systemctl", "status", "hezproxy")
	cmd.Stdout = &stdOut

	err := cmd.Run()

	if err != nil {
		return err
	}

	fmt.Printf(stdOut.String())

	return nil
}

func StatusCommand() *cobra.Command {
	return core.CreateCommand(&cobra.Command{
		Use:   "run",
		Short: "Returns the status of the reverse proxy.",
		Long:  `Returns the status of the reverse proxy.`,
	}, nil, runStatusCommand)
}
