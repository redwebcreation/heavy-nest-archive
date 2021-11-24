package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wormable/nest/ansi"
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
	ansi.Check(err)

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
	executable, err := os.Executable()
	ansi.Check(err)

	usingGoRun := strings.HasPrefix(executable, "/tmp/go-build")

	if usingGoRun {
		fmt.Printf("%sYou're running this command using go run, you won't see any effects.\n%s", ansi.Yellow.Fg(), ansi.Reset)
		fmt.Printf("%sPlease build nest's binary and then run it.\n\n%s", ansi.Yellow.Fg(), ansi.Reset)
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
		fmt.Printf("%sNest is already up to date.%s", ansi.Green.Fg(), ansi.Reset)
		return nil
	}

	binary := latestRelease.Assets[0]

	if binary.State != "uploaded" {
		fmt.Println("The binary for the latest release is still being uploaded.")
		fmt.Println("Please try again in a few seconds.")
	}

	fmt.Printf("Downloading %s.\n", binary.BrowserDownloadUrl)

	response, err := http.Get(binary.BrowserDownloadUrl)
	ansi.Check(err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	ansi.Check(err)

	err = os.WriteFile(executable+"_updated", body, os.FileMode(0777))
	ansi.Check(err)

	err = os.Rename(executable+"_updated", executable)
	ansi.Check(err)

	_ = os.Chmod(executable, os.FileMode(0777))

	fmt.Printf("%sSuccessfully updated nest to %s\n", ansi.Green.Fg(), latestRelease.TagName+ansi.Reset)

	return nil
}

func SelfUpdateCommand() *cobra.Command {
	return Decorate(&cobra.Command{
		Use:   "self-update [version]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "update nest to the latest version",
	}, runSelfUpdateCommand, func(c *cobra.Command) {
		c.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Force the update")
	})
}
