package cli

import (
	"errors"
	"fmt"
	"github.com/redwebcreation/hez/ansi"
	"github.com/redwebcreation/hez/core"
	"github.com/redwebcreation/hez/globals"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var dryRun bool
var force bool
var draft bool
var prerelease bool

func RunSelfUpdateCommand(_ *cobra.Command, args []string) error {
	err := core.ElevateProcess()

	if err != nil {
		return err
	}

	if dryRun {
		ansi.Warning("Dry running the command, nothing will be executed.")
	}

	executable, _ := os.Executable()
	usingGoRun := strings.HasPrefix(executable, "/tmp/go-build")

	version := ""

	if len(args) > 0 {
		version = args[0]
	}

	if usingGoRun {
		ansi.Warning("You're running this command using go run, you won't see any effects.")
		ansi.Warning("Please build the binary and then run it.")
		return nil
	}

	//goland:noinspection GoBoolExpressions
	if globals.Version == "(development)" {
		ansi.Warning("You're using the development build of Hez.")
		ansi.Warning("Please, specify a version to use when building the binary.")
		ansi.Warning("go build -ldflags=\"-X github.com/redwebcreation/hez/globals.Version=$(git describe --tags)\"")
		return nil
	}

	releases, err := core.Repository.Releases(core.ReleaseFilter{
		Draft:      draft,
		Prerelease: prerelease,
		Version:    version,
	})

	if err != nil {
		return err
	}

	latestRelease := releases[0]

	if !force && latestRelease.TagName == globals.Version {
		fmt.Println("You're already using the latest version.")
		return nil
	}

	binary := latestRelease.Assets[0]

	if binary.State != "uploaded" {
		return errors.New("The binary for the latest release is still being uploaded. \nPlease try again in a few seconds.")
	}

	fmt.Printf("Downloading %s.\n", binary.BrowserDownloadUrl)

	response, err := http.Get(binary.BrowserDownloadUrl)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if dryRun {
		ansi.Success("Successfully updated Hez to the version " + latestRelease.TagName)
		return nil
	}

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

	ansi.Success("Successfully updated Hez to the version " + latestRelease.TagName)

	return nil
}

func SelfUpdateCommand() *cobra.Command {
	command := core.CreateCommand(&cobra.Command{
		Aliases: []string{"selfupdate", "update"},
		Use:     "self-update [version]",
		Short:   "Updates Hez to the latest version.",
		Long:    `Updates Hez to the latest version or the one given as the first argument.`,
	}, func(command *cobra.Command) {
		command.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run Hez's update process.")
		command.Flags().BoolVarP(&force, "force", "f", false, "Force the update.")
		command.Flags().BoolVarP(&draft, "edge", "e", false, "Updates to the main branch build.")
		command.Flags().BoolVarP(&prerelease, "prerelease", "p", false, "Updates to the latest prerelease.")
	}, RunSelfUpdateCommand)

	return command
}
