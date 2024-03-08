package providers

import "github.com/nerzhul/releasetracker/pkg/release"

type ReleaseProvider interface {
	GetReleases(group string, repo string, maxReleases int) (release.ReleaseList, error)
}
