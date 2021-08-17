package internal

import (
	"bytes"
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

func EnableProxy() error {
	supervisorConfig, err := GetSupervisordConfig()

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

	err = RunCommand("systemctl", "daemon-reload")
	if err != nil {
		return err
	}
	err = RunCommand("systemctl", "enable", "hezproxy")
	if err != nil {
		return err
	}
	return RunCommand("service", "hezproxy", "start")
}

func DisableProxy() error {
	configName := "/etc/systemd/system/hezproxy.service"

	err := RunCommand("systemctl", "stop", "hezproxy")
	if err != nil {
		return err
	}
	err = RunCommand("systemctl", "disable", "hezproxy")
	if err != nil {
		return err
	}
	err = RunCommand("systemctl", "daemon-reload")
	if err != nil {
		return err
	}

	return os.Remove(configName)
}

func GetSupervisordConfig() (string, error) {
	stub := `[Unit]
Description=Hez Proxy Server
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=[user]
ExecStart=[executable] proxy run

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

	return stub, nil
}

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	return cmd.Run()
}
