package release

type Release struct {
	Version string
	ReleaseDate string
}

type ReleaseList struct {
	Provider string
	Group string
	Repo string
	Releases []Release
}