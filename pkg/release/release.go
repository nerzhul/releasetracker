package release

type Release struct {
	Version     string
	ReleaseDate string
}

type ReleaseMetadata struct {
	Provider string
	Group    string
	Repo     string
}

type ReleaseList struct {
	Metadata ReleaseMetadata
	Releases []Release
}
