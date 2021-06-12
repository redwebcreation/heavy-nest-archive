package cli

import (
	"encoding/json"
	"errors"
	"github.com/redwebcreation/hez2/core"
	"github.com/redwebcreation/hez2/globals"
	"github.com/redwebcreation/hez2/util"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var dryRun bool
var force bool
var draft bool
var prerelease bool

type Repository string

type Asset struct {
	Url                string    `json:"url"`
	Id                 int       `json:"id"`
	NodeId             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	Uploader           User      `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadUrl string    `json:"browser_download_url"`
}

type User struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Release struct {
	Url             string    `json:"url"`
	AssetsUrl       string    `json:"assets_url"`
	UploadUrl       string    `json:"upload_url"`
	HtmlUrl         string    `json:"html_url"`
	Id              int       `json:"id"`
	Author          User      `json:"author"`
	NodeId          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []Asset   `json:"assets"`
	TarballUrl      string    `json:"tarball_url"`
	ZipballUrl      string    `json:"zipball_url"`
	Body            string    `json:"body"`
}

type ReleaseFilter struct {
	Draft      bool
	Prerelease bool
	Version    string
}

func NewGithubRequest(url string, data interface{}) error {
	response, err := http.Get("https://api.github.com/" + url)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return err
	}

	return nil
}

func (repository Repository) Releases(filter ReleaseFilter) ([]Release, error) {
	var releases []Release

	err := NewGithubRequest("repos/"+string(repository)+"/releases", &releases)

	if err != nil {
		return releases, err
	}

	var filtered []Release

	for _, release := range releases {
		if release.TagName != filter.Version && filter.Version != "" {
			continue
		}

		if release.Draft != filter.Draft {
			continue
		}

		if release.Prerelease != filter.Prerelease {
			continue
		}

		filtered = append(filtered, release)
	}

	return filtered, nil
}

var Repo = Repository("redwebcreation/hez")

func RunSelfUpdateCommand(_ *cobra.Command, args []string) error {
	if dryRun {
		core.Ansi.Warning("Dry running the command, nothing will be executed.")
	}
	executable, _ := os.Executable()
	usingGoRun := strings.HasPrefix(executable, "/tmp/go-build")

	version := ""

	if len(args) > 0 {
		version = args[0]
	}

	if usingGoRun {
		core.Ansi.Warning("You're running this command using go run, you won't see any effects.")
		core.Ansi.Warning("Please build the binary and then run it.")
		return nil
	}

	//goland:noinspection GoBoolExpressions
	if globals.Version == "(development)" {
		core.Ansi.Warning("You're using the development build of Hez.")
		core.Ansi.Warning("Please, specify a version to use when building the binary.")
		core.Ansi.Warning("go build -ldflags=\"-X github.com/redwebcreation/hez2/core.Version=$(git describe --tags)\"")
		return nil
	}

	releases, err := Repo.Releases(ReleaseFilter{
		Draft:      draft,
		Prerelease: prerelease,
		Version:    version,
	})

	if err != nil {
		return err
	}

	latestRelease := releases[0]

	if !force && latestRelease.TagName == globals.Version {
		core.Ansi.Success("You're already using Hez " + latestRelease.TagName + ".")
		return nil
	}

	binary := latestRelease.Assets[0]

	if binary.State != "uploaded" {
		return errors.New("The binary for the latest release is still being uploaded. \nPlease try again in a few seconds.")
	}

	core.Ansi.Printf("Downloading %s.\n", binary.BrowserDownloadUrl)

	response, err := http.Get(binary.BrowserDownloadUrl)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if dryRun {
		core.Ansi.Success("Successfully updated Hez to the version " + latestRelease.TagName)
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

	core.Ansi.Success("Successfully updated Hez to the version " + latestRelease.TagName)

	return nil
}

func SelfUpdateCommand() *cobra.Command {
	command := util.CreateCommand(&cobra.Command{
		Use:   "self-update [version]",
		Short: "Updates Hez to the latest version.",
		Long:  `Updates Hez to the latest version or the one given as the first argument.`,
	}, func(command *cobra.Command) {
		command.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run Hez's update process.")
		command.Flags().BoolVarP(&force, "force", "f", false, "Force the update.")
		command.Flags().BoolVarP(&draft, "edge", "e", false, "Updates to the main branch build.")
		command.Flags().BoolVarP(&prerelease, "prerelease", "p", false, "Updates to the latest prerelease.")
	}, RunSelfUpdateCommand)

	return command
}
