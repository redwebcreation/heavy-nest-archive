package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd/ui"
	"github.com/wormable/nest/globals"
)

var forceUpdate bool

type githubRepository string

type release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		State              string `json:"state"`
		BrowserDownloadUrl string `json:"browser_download_url"`
	}
}

const nestRepository = githubRepository("wormable/nest")

func (g githubRepository) newRequest(url string, data interface{}) error {
	response, err := http.Get("https://api.github.com/" + url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		var c interface{}
		_ = json.Unmarshal(body, &c)
		pp, _ := json.MarshalIndent(c, "", "  ")
		fmt.Println(string(pp))
	}

	return err
}

func (g githubRepository) Releases(version string) []release {
	var releases []release
	err := g.newRequest("repos/"+g.String()+"/releases", &releases)
	ui.Check(err)

	var filtered []release

	for _, release := range releases {
		if release.TagName != version && version != "" {
			continue
		}

		filtered = append(filtered, release)
	}

	return filtered
}

func (g githubRepository) String() string {
	return string(g)
}

func runSelfUpdateCommand(_ *cobra.Command, args []string) error {
	ElevateProcess()

	executable, err := os.Executable()
	ui.Check(err)

	usingGoRun := strings.HasPrefix(executable, "/tmp/go-build")

	if usingGoRun {
		fmt.Printf("%sYou're running this command using go run, you won't see any effects.\n%s", ui.Yellow.Fg(), ui.Stop)
		fmt.Printf("%sPlease build nest's binary and then run it.\n\n%s", ui.Yellow.Fg(), ui.Stop)
	}

	var versionNeeded string

	if len(args) > 0 {
		versionNeeded = args[0]
	}

	releases := nestRepository.Releases(versionNeeded)
	latestRelease := releases[0]

	if len(releases) == 0 {
		return fmt.Errorf("no releases found using the given filters")
	}

	if !forceUpdate && latestRelease.TagName == globals.Version {
		fmt.Printf("%sNest is already up to date.%s", ui.Green.Fg(), ui.Stop)
		return nil
	}

	binary := latestRelease.Assets[0]

	if binary.State != "uploaded" {
		fmt.Println("The binary for the latest release is still being uploaded.")
		fmt.Println("Please try again in a few seconds.")
	}

	fmt.Printf("Downloading %s.\n", binary.BrowserDownloadUrl)

	response, err := http.Get(binary.BrowserDownloadUrl)
	ui.Check(err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	ui.Check(err)

	err = os.WriteFile(executable+"_updated", body, os.FileMode(0777))
	ui.Check(err)

	err = os.Rename(executable+"_updated", executable)
	ui.Check(err)

	_ = os.Chmod(executable, os.FileMode(0777))

	fmt.Printf("%sSuccessfully updated nest to %s%s\n", ui.Green.Fg(), latestRelease.TagName, ui.Stop)

	return nil
}

func SelfUpdateCommand() *cobra.Command {
	return CreateCommand(&cobra.Command{
		Use:   "self-update [version]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Update nest to the latest version",
	}, func(c *cobra.Command) {
		c.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Force the update")
	}, runSelfUpdateCommand)
}
