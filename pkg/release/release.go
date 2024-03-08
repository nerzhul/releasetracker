package release

type Release struct {
	Version string
	ReleaseDate string
}

type ReleaseList struct {
	Group string
	Repo string
	Releases []Release
}