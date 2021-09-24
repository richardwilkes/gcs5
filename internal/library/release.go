package library

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// Release holds information about a single release of a GitHub repo.
type Release struct {
	Version     Version
	Notes       string
	ZipFileURL  string
	CheckFailed bool
}

// HasUpdate returns true if there is an update available.
func (r *Release) HasUpdate() bool {
	return !r.CheckFailed && r.Version != Version{}
}

// LoadReleases loads the list of releases available from a given GitHub repo.
func LoadReleases(ctx context.Context, client *http.Client, githubAccountName, repoName string, currentVersion Version, filter func(version Version, notes string) bool) ([]Release, error) {
	if githubAccountName == "*" {
		return nil, nil
	}
	var versions []Release
	uri := "https://api.github.com/repos/" + githubAccountName + "/" + repoName + "/releases"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, errs.NewWithCause("unable to create GitHub API request "+uri, err)
	}
	var rsp *http.Response
	if rsp, err = client.Do(req); err != nil {
		return nil, errs.NewWithCause("GitHub API request failed "+uri, err)
	}
	defer xio.DiscardAndCloseIgnoringErrors(rsp.Body)
	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		return nil, errs.New("unexpected response code from GitHub API " + uri + " -> " + rsp.Status)
	}
	var releases []struct {
		TagName    string `json:"tag_name"`
		Body       string `json:"body"`
		ZipBallURL string `json:"zipball_url"`
	}
	if err = json.NewDecoder(rsp.Body).Decode(&releases); err != nil {
		return nil, errs.NewWithCause("unable to decode response from GitHub API "+uri, err)
	}
	for _, one := range releases {
		if strings.HasPrefix(one.TagName, "v") {
			if version := VersionFromString(one.TagName[1:]); version != (Version{}) && (currentVersion == version || currentVersion.Less(version)) {
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
	if len(versions) > 1 && versions[len(versions)-1].Version == currentVersion {
		versions = versions[:len(versions)-1]
	}
	return versions, nil
}
