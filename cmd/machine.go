package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/ansi"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func runMachine(_ *cobra.Command, _ []string) error {
	// This is complex and asynchronous because in the future we may need to make many requests
	// to a given backend.
	var wg sync.WaitGroup
	done := func() { wg.Done() }
	details := []func(map[string]string, func()){
		getPublicIP,
		getMemoryLimit,
		getAvailableProcessors,
	}

	detailsLength := len(details)

	data := make(map[string]string, detailsLength)

	for _, detail := range details {
		go detail(data, done)
	}

	wg.Add(detailsLength)
	wg.Wait()

	for k, v := range data {
		if v == "error" {
			fmt.Printf("%s: %s\n", ansi.Gray.Fg()+k, ansi.Red.Fg()+"errored"+ansi.Reset)
			continue
		}
		fmt.Printf("%s: %s\n", ansi.Gray.Fg()+k, ansi.White.Fg()+strings.TrimSpace(v)+ansi.Reset)
	}

	return nil
}

func getMemoryLimit(data map[string]string, done func()) {
	defer done()
	// returns the memory limit in kilobytes
	rawLimit, err := exec.Command("sh", "-c", "cat /proc/meminfo | grep MemTotal | awk '{print $2}'").CombinedOutput()
	if err != nil {
		data["memory limit"] = "error"
		return
	}

	limit, err := strconv.ParseFloat(strings.TrimSpace(string(rawLimit)), 64)
	if err != nil {
		data["memory limit"] = "error"
		return
	}

	data["memory limit"] = fmt.Sprintf("%.0fgb (%.f)", limit/(1024.0*1024.0), limit)
}

func getAvailableProcessors(data map[string]string, done func()) {
	defer done()
	out, err := exec.Command("nproc").CombinedOutput()
	if err != nil {
		data["available processors"] = "error"
		return
	}

	data["available processors"] = string(out)
}

func MachineCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "machine",
		Short: "Get details about the current machine",
	}, runMachine, nil)
}

func getPublicIP(data map[string]string, done func()) {
	defer done()

	response, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		data["public ip"] = "error"
	}

	defer response.Body.Close()

	rawPublicIp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		data["public ip"] = "error"
		return
	}

	data["public ip"] = string(rawPublicIp)
}
