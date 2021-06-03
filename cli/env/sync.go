package env

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"os"
)

type stubFile struct {
	Path      string
	Directory bool
}

func runSyncCommand(_ *cobra.Command, _ []string) {
	config := core.GetConfig()
	validEnvironments := make([]string, len(config.Applications))

	for _, application := range config.Applications {
		appEnv := application.Env
		validEnvironments = append(validEnvironments, core.EnvironmentPath(appEnv))

		fmt.Println("[" + appEnv + "]")

		if core.EnvironmentExists(appEnv) {
			fmt.Println("  - already synchronized.")
			continue
		}

		stubs := [5]stubFile{
			{
				Path:      core.EnvironmentPath(appEnv),
				Directory: true,
			},
			{
				Path:      core.EnvironmentPath(appEnv) + "/current",
				Directory: true,
			},
			{
				Path:      core.EnvironmentPath(appEnv) + "/staging",
				Directory: true,
			},
			{
				Path:      core.EnvironmentPath(appEnv) + "/current/.env",
				Directory: false,
			},
			{
				Path:      core.EnvironmentPath(appEnv) + "/staging/.env",
				Directory: false,
			},
		}

		for _, stub := range stubs {
			if stub.Directory {
				err := os.Mkdir(stub.Path, os.FileMode(0777))

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("  - new directory " + stub.Path)
				}
			} else {
				err := os.WriteFile(stub.Path, nil, os.FileMode(0777))

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("  - new " + stub.Path)
				}
			}
		}
	}

	files, err := os.ReadDir(core.EnvironmentPath(""))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("[__overflow__]")
	hasOverflow := false
	for _, file := range files {
		path := core.EnvironmentPath(file.Name())
		if !file.IsDir() {
			hasOverflow = true
			_ = os.Remove(path)
			fmt.Println("  - removed " + path)
			continue
		}

		if !contains(validEnvironments, path) {
			hasOverflow = true
			_ = os.RemoveAll(path)
			fmt.Println("  - removed " + path)
		}
	}

	if !hasOverflow {
		fmt.Println("  - no overflow")
	}

	fmt.Println("Everything is in sync.")
}

func initSyncCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Deletes unused environments. Creates new ones.",
		Run:   runSyncCommand,
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
