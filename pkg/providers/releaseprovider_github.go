package providers

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/nerzhul/releasetracker/pkg/release"
)

type GithubReleaseProvider struct{
	client *github.Client
}

func NewGithubReleaseProvider() *GithubReleaseProvider {
	return &GithubReleaseProvider{
		client: github.NewClient(nil),
	}
}

func (p *GithubReleaseProvider) GetReleases(group string, repo string, maxReleases int) (*release.ReleaseList, error) {
	ghReleases, resp, err := p.client.Repositories.ListReleases(context.TODO(), group, repo, &github.ListOptions{
		PerPage: maxReleases,
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch releases from github (%s/%s): %s", group, repo, resp.Status)
	}

	releasesList := &release.ReleaseList{
		Group: group,
		Repo: repo,
		Releases: make([]release.Release, 0),
	}

	for _, ghRelease := range ghReleases {
		// No tag name, no release date, or prerelease
		if ghRelease.TagName == nil || ghRelease.PublishedAt == nil ||
			ghRelease.Prerelease != nil && *ghRelease.Prerelease {
			continue
		}

		releasesList.Releases = append(releasesList.Releases, release.Release{
			Version: *ghRelease.TagName,
			ReleaseDate: ghRelease.PublishedAt.Format("2006-01-02"),
		})
	}
	return releasesList, nil
}

