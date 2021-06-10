package cli

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"strconv"
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

	fmt.Println("Total memory: " + FormatBytes(ram))
	fmt.Println("Free memory: " + FormatBytes(sysinfo.Freeram))
	return nil
}

func InfoCommand() *cobra.Command {
	command := CreateCommand(&cobra.Command{
		Use:   "info [node name]",
		Short: "Displays various metrics about your system.",
		Long:  `Display various metrics about the master's hardware such as available memory, cpu cores...`,
	}, nil, RunInfoCommand)

	return command
}

func FormatBytes(bytes uint64) string {
	var format int64 = 1

	for float64(bytes) > math.Pow(float64(format), float64(10)) {
		format += 1
	}

	inBytes := math.Round(float64(bytes / uint64(format)))
	asString := strconv.FormatUint(uint64(inBytes/uint64(format)), 10)
	switch format - 2 {
	case 9:
		asString += "Go"
	}

	return asString
}

func GetSysInfo() (*syscall.Sysinfo_t, error) {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)

	if err != nil {
		return nil, err
	}

	return in, nil
}
