package core

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	systemUser "os/user"
	"strings"
)

func IsProxyEnabled() (bool, error) {
	var stdOut bytes.Buffer
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service")
	cmd.Stdout = &stdOut
	err := cmd.Run()

	lines := strings.Split(strings.TrimSpace(stdOut.String()), "\n")

	if err != nil {
		return false, err
	}

	for _, line := range lines {
		if strings.Contains(line, "hezproxy.service") && strings.Contains(line, "enabled") {
			return true, nil
		}
	}

	return false, nil
}

func EnableProxy(proxy Proxy) error {
	supervisorConfig, err := GetSupervisordConfig(proxy)

	if err != nil {
		return err
	}

	configName := "/etc/systemd/system/hezproxy.service"
	configFile, err := os.Stat(configName)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if configFile != nil {
		err := os.Remove(configName)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(configName, []byte(supervisorConfig), os.FileMode(0744))

	if err != nil {
		return err
	}

	var commands = []string{
		"systemctl daemon-reload",
		"systemctl enable hezproxy",
		"service hezproxy start",
	}

	for _, command := range commands {
		err := runCommand(command)

		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(command string) error {
	var stdErr bytes.Buffer
	commandParts := strings.Split(command, " ")
	name := commandParts[0]
	args := commandParts[1:]

	cmd := exec.Command(name, args...)
	cmd.Stderr = &stdErr

	err := cmd.Run()

	if stdErr.Len() > 0 {
		return errors.New(stdErr.String())
	}

	return err
}

func DisableProxy(proxy Proxy) error {
	configName := "/etc/systemd/system/hezproxy.service"

	var commands = []string{
		"systemctl stop hezproxy",
		"systemctl disable hezproxy",
		"systemctl daemon-reload",
	}

	for _, command := range commands {
		err := runCommand(command)

		if err != nil {
			return err
		}
	}

	err := os.Remove(configName)

	return err
}

func GetSupervisordConfig(proxy Proxy) (string, error) {
	stub := `[Unit]
Description=Hez Proxy Server
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=[user]
ExecStart=[executable] proxy run --port [port] --ssl [ssl] [selfSigned]

[Install]
WantedBy=multi-user.target`

	executable, err := os.Executable()

	if err != nil {
		return "", err
	}

	user, err := systemUser.Current()

	if err != nil {
		return "", err
	}

	stub = strings.Replace(stub, "[user]", user.Username, 1)
	stub = strings.Replace(stub, "[executable]", executable, 1)
	stub = strings.Replace(stub, "[port]", proxy.Port, 1)
	stub = strings.Replace(stub, "[ssl]", proxy.Ssl, 1)

	if proxy.SelfSigned {
		stub = strings.Replace(stub, "[selfSigned]", "--self-signed", 1)
	} else {
		stub = strings.Replace(stub, "[selfSigned]", "", 1)
	}

	return stub, nil
}

func IsRunningAsRoot() bool {
	return os.Geteuid() == 0
}
