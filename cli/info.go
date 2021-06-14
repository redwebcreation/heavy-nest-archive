package cli

import (
	"errors"
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func RunInfoCommand(_ *cobra.Command, args []string) error {
	var node string
	if len(args) > 0 {
		node = args[0]
	} else {
		node = "master"
	}

	if node != "master" {
		return errors.New("node not found")
	}

	sysinfo, err := GetSysInfo()
	if err != nil {
		return err
	}

	ram := sysinfo.Totalram * uint64(sysinfo.Unit)

	externalIp, err := GetExternalIp()

	if err != nil {
		return err
	}

	fmt.Println("external_ip: ", externalIp)
	fmt.Println("total_memory: " + FormatBytes(ram))
	fmt.Println("free_memory: " + FormatBytes(sysinfo.Freeram))

	cpu, err := GetCpu()

	if err != nil {
		return err
	}

	fmt.Println("processor: " + cpu.ModelName)
	fmt.Println("processor_cores: " + strconv.FormatUint(cpu.Cores, 10))
	return nil
}

func InfoCommand() *cobra.Command {
	command := core.CreateCommand(&cobra.Command{
		Use:   "info [node name]",
		Short: "Displays various metrics about your system",
		Long:  `Display various metrics about the master's hardware such as available memory, cpu cores...`,
	}, nil, RunInfoCommand)

	return command
}

func FormatBytes(b uint64) string {
	const unit = 1000

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0

	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func GetSysInfo() (*syscall.Sysinfo_t, error) {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)

	if err != nil {
		return nil, err
	}

	return in, nil
}

type Core struct {
	ModelName string
	Cores     uint64
}

func GetCpu() (Core, error) {
	contents, err := os.ReadFile("/proc/cpuinfo")

	if err != nil {
		return Core{}, err
	}

	core := Core{}

	for _, line := range strings.Split(string(contents), "\n") {
		if line == "" {
			break
		}

		keyValue := strings.Split(line, ": ")
		key := strings.TrimSpace(keyValue[0])
		value := keyValue[1]

		if key == "cpu cores" {
			cpuCures, _ := strconv.ParseUint(value, 10, 64)

			core.Cores = cpuCures
		}

		if key == "model name" {
			core.ModelName = value
		}
	}

	return core, nil
}

func GetExternalIp() (string, error) {
	// TODO: MAYBE Implement a Voter/Consensus system
	response, err := http.Get("http://checkip.amazonaws.com/")
	if err != nil {
		return "", err

	}

	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
