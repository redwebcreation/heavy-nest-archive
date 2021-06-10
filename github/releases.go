package github

type Asset struct {
	Url                string
	Id                 int
	NodeId             string
	Name               string
	Label              string
	Uploader           User
	ContentType        string
	State              string
	Size               int
	DownloadCount      int
	CreatedAt          string
	UpdatedAt          string
	BrowserDownloadUrl string
}

type User struct {
	Login             string
	Id                int
	NodeId            string
	AvatarUrl         string
	GravatarId        string
	Url               string
	HtmlUrl           string
	FollowersUrl      string
	FollowingUrl      string
	GistsUrl          string
	StarredUrl        string
	SubscriptionsUrl  string
	OrganizationsUrl  string
	ReposUrl          string
	EventsUrl         string
	ReceivedEventsUrl string
	Type              string
	SiteAdmin         bool
}

type Release struct {
	Url             string
	AssetsUrl       string
	UploadUrl       string
	HtmlUrl         string
	Id              int
	Author          User
	NodeId          string
	TagName         string
	TargetCommitish string
	Name            string
	Draft           bool
	Prerelease      bool
	CreatedAt       string
	PublishedAt     string
	Assets          []Asset
	TarballUrl      string
	ZipballUrl      string
	Body            string
}

type ReleaseFilter struct {
	Draft      bool
	Prerelease bool
	Version    string
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
