package cli

import (
	"errors"
	"fmt"
	"github.com/redwebcreation/hez/globals"
	"github.com/redwebcreation/hez/internal"
	ansi2 "github.com/redwebcreation/hez/internal/ui"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var force bool
var edge bool
var prerelease bool

func RunSelfUpdateCommand(_ *cobra.Command, args []string) error {
	err := internal.ElevateProcess()

	if err != nil {
		return err
	}

	executable, _ := os.Executable()
	usingGoRun := strings.HasPrefix(executable, "/tmp/go-build")

	version := ""

	if len(args) > 0 {
		version = args[0]
	}

	if usingGoRun {
		ansi2.Warning("You're running this command using go run, you won't see any effects.")
		ansi2.Warning("Please build the binary and then run it.")
		return nil
	}

	var response *http.Response
	var latestRelease internal.Release

	if edge {
		response, err = UpdateToEdge()

		if err != nil {
			return err
		}
	} else {
		releases, err := internal.Repository.Releases(internal.ReleaseFilter{
			Prerelease: prerelease,
			Version:    version,
		})

		if err != nil {
			return err
		}

		if len(releases) == 0 {
			return errors.New("no releases found using the given filters")
		}

		latestRelease = releases[0]

		if !force && latestRelease.TagName == globals.Version {
			fmt.Println("You're already using the latest version.")
			return nil
		}

		binary := latestRelease.Assets[0]

		if binary.State != "uploaded" {
			return errors.New("The binary for the latest release is still being uploaded. \nPlease try again in a few seconds.")
		}

		fmt.Printf("Downloading %s.\n", binary.BrowserDownloadUrl)

		response, err = http.Get(binary.BrowserDownloadUrl)

		if err != nil {
			return err
		}
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = os.WriteFile(executable+"_updated", body, os.FileMode(0777))

	if err != nil {
		return err
	}

	err = os.Rename(executable+"_updated", executable)

	if err != nil {
		return err
	}

	_ = os.Chmod(executable, os.FileMode(0777))

	var updatedVersion string
	if edge {
		updatedVersion = "development version"
	} else {
		updatedVersion = "version " + latestRelease.TagName
	}

	ansi2.Success("Successfully updated Hez to the " + updatedVersion)

	return nil
}

func UpdateToEdge() (*http.Response, error) {
	url := "https://raw.githubusercontent.com/" + string(internal.Repository) + "/master/hez"

	fmt.Printf("Downloading %s.\n", url)

	return http.Get(url)
}

func SelfUpdateCommand() *cobra.Command {
	command := internal.CreateCommand(&cobra.Command{
		Aliases: []string{"selfupdate", "update"},
		Use:     "self-update [version]",
		Short:   "Updates Hez to the latest version.",
		Long:    `Updates Hez to the latest version or the one given as the first argument.`,
	}, func(command *cobra.Command) {
		command.Flags().BoolVarP(&force, "force", "f", false, "Force the update.")
		command.Flags().BoolVarP(&edge, "edge", "e", false, "Updates to the main branch build.")
		command.Flags().BoolVarP(&prerelease, "prerelease", "p", false, "Updates to the latest prerelease.")
	}, RunSelfUpdateCommand)

	return command
}
