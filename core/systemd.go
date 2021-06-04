package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	systemUser "os/user"
	"strings"
)

func IsProxyEnabled() bool {
	var stdOut bytes.Buffer
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service")
	cmd.Stdout = &stdOut
	err := cmd.Run()

	lines := strings.Split(strings.TrimSpace(stdOut.String()), "\n")

	if err != nil {
		return false
	}

	for _, line := range lines {
		if strings.Contains(line, "hezproxy.service") && strings.Contains(line, "enabled") {
			return true
		}
	}

	return false
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
		fmt.Println("Running [" + command + "]")

		if err != nil {
			fmt.Println("ERROR HERE")
			return err
		}
	}

	return nil
}

func runCommand(command string) error {
	commandParts := strings.Split(command, " ")
	name := commandParts[0]
	args := commandParts[1:]

	cmd := exec.Command(name, args...)

	err := cmd.Run()

	return err
}

func DisableProxy() error {
	configName := "/etc/systemd/system/hezproxy.service"

	var commands = []string{
		"systemctl stop hezproxy",
		"systemctl disable hezproxy",
		"systemctl daemon-reload",
	}

	for _, command := range commands {
		err := runCommand(command)
		fmt.Println("Running [" + command + "]")

		if err != nil {
			fmt.Println("ERROR HERE")
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

	if *proxy.SelfSigned {
		stub = strings.Replace(stub, "[selfSigned]", "--self-signed", 1)
	} else {
		stub = strings.Replace(stub, "[selfSigned]", "", 1)
	}

	return stub, nil
}

func IsRunningAsRoot() bool {
	return os.Geteuid() == 0
}
