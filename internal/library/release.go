package library

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
)

// Release holds information about a single release of a GitHub repo.
type Release struct {
	Version            Version
	Notes              string
	ZipFileURL         string
	UnableToAccessRepo bool
}

// LoadReleases loads the list of releases available from a given GitHub repo.
func LoadReleases(ctx context.Context, client *http.Client, githubAccountName, repoName string, currentVersion Version, filter func(version Version, notes string) bool) []Release {
	var versions []Release
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/"+githubAccountName+"/"+repoName+"/releases", nil)
	if err != nil {
		jot.Error(errs.NewWithCause("unable to create GitHub API request", err))
		return nil
	}
	var rsp *http.Response
	if rsp, err = client.Do(req); err != nil {
		jot.Error(errs.NewWithCause("GitHub API request failed", err))
		return nil
	}
	defer xio.DiscardAndCloseIgnoringErrors(rsp.Body)
	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		jot.Error(errs.New("unexpected response code from GitHub API: " + rsp.Status))
		return nil
	}
	var releases []struct {
		TagName    string `json:"tag_name"`
		Body       string `json:"body"`
		ZipBallURL string `json:"zipball_url"`
	}
	if err = json.NewDecoder(rsp.Body).Decode(&releases); err != nil {
		jot.Error(errs.NewWithCause("unable to decode response from GitHub API", err))
		return nil
	}
	for _, one := range releases {
		if strings.HasPrefix(one.TagName, "v") {
			if version := VersionFromString(one.TagName[1:]); version != (Version{}) && currentVersion.Less(version) {
				if filter == nil || !filter(version, one.Body) {
					versions = append(versions, Release{
						Version:    version,
						Notes:      one.Body,
						ZipFileURL: one.ZipBallURL,
					})
				}
			}
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[j].Version.Less(versions[i].Version)
	})
	return versions
}

// DistillReleases distills the release list down down a single representative release.
func DistillReleases(releases []Release) Release {
	switch len(releases) {
	case 0:
		return Release{UnableToAccessRepo: true}
	case 1:
		return releases[0]
	default:
		release := releases[0]
		for _, one := range releases[1:] {
			release.Notes += "\n\n## Version " + one.Version.String() + "\n" + one.Notes
		}
		return release
	}
}
