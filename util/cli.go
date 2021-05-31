package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Command(arguments ...string) (string, error) {
	var out bytes.Buffer
	var errBuffer bytes.Buffer
	cmd := exec.Command(arguments[0], arguments[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &errBuffer
	err := cmd.Run()

	fmt.Println("Buffer", errBuffer.String())

	if err != nil {
		return "", err
	}


	return out.String(), nil
}
