package core

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func ElevateProcess() error {
	cmd := exec.Command("sudo", "touch", "/tmp/upgrade-process")

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunCommand(command ...string) error {
	name := command[0]
	args := command[1:]

	var stderr bytes.Buffer

	cmd := exec.Command(name, args...)

	cmd.Stderr = &stderr

	err := cmd.Run()

	if stderr.Len() > 0 {
		return errors.New(stderr.String())
	}

	return err
}
