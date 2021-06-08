package update

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/hez/core"
	box "github.com/redwebcreation/hez/core/embed"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Release struct {
	Version    string `json:"tag_name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		Name  string `json:"name"`
		Url   string `json:"browser_download_url"`
		State string `json:"state"`
		Size  int    `json:"size"`
	}
}

var draft bool
var prerelease bool

func run(_ *cobra.Command, _ []string) {
	currentVersion := strings.TrimSpace(string(box.Get("/version")))
	latestRelease := FindRelease(draft, prerelease)

	if latestRelease.Version == currentVersion {
		fmt.Printf("You are using the latest version of Hez (%s).\n", latestRelease.Version)
		return
	}

	fmt.Printf("Updating to the latest version %s.\n", latestRelease.Version)

	hezBinary := latestRelease.Assets[0]

	if hezBinary.State != "uploaded" {
		fmt.Println("The assets for this release are still being uploaded.\n Please try again in a minute.")
		return
	}

	fmt.Printf("Downloading %s.\n", hezBinary.Url)

	response, err := http.Get(hezBinary.Url)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	currentExecutable, _ := os.Executable()

	err = core.CreateFile("updated_hez", body, os.FileMode(0777))

	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.Rename(core.GetPathFor("updated_hez"), currentExecutable)

	if err != nil {
		fmt.Println(err)
		return
	}

	_ = os.Chmod(currentExecutable, os.FileMode(0777))

	fmt.Println("You are now using the latest version of Hez.")
}

func FindRelease(draft bool, prerelease bool) Release {
	releases := GetReleases()
	var fittingRelease Release

	for _, release := range releases {
		if release.Draft == draft && release.Prerelease == prerelease {
			fittingRelease = release
			break
		}
	}
	return fittingRelease
}

func NewCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates Hez to the latest version.",
		Run:   run,
	}

	updateCmd.Flags().BoolVarP(&draft, "draft", "d", false, "Update to the latest draft release")
	updateCmd.Flags().BoolVarP(&prerelease, "prerelease", "p", false, "Update to the latest prerelease")

	return updateCmd
}

func GetReleases() []Release {
	var releases []Release

	response, err := http.Get("https://api.github.com/repos/redwebcreation/hez/releases")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(body, &releases)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return releases
}
