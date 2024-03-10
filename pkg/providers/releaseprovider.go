package providers

import "github.com/nerzhul/releasetracker/pkg/release"

type ReleaseProvider interface {
	GetReleases(provider string, group string, repo string, maxReleases int) (release.ReleaseList, error)
	RecordReleases(provider string, group string, repo string, releases *release.ReleaseList) error
}
