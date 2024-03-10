package releasesyncer

import "github.com/nerzhul/releasetracker/pkg/providers"

type ReleaseSyncer struct {
	db *ReleaseSyncerDB
	ghProvider *providers.GithubReleaseProvider
}

func NewReleaseSyncer(databaseURL string) *ReleaseSyncer {
	return &ReleaseSyncer{
		db: NewReleaseSyncerDB(databaseURL),
		ghProvider: providers.NewGithubReleaseProvider(),
	}
}

func (s *ReleaseSyncer) SyncReleases() error {
	subs, err := s.db.getSubscribedReleases(providers.GithubReleaseProviderName)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		releases, err := s.ghProvider.GetReleases("", sub.Group, sub.Repo, 50)
		if err != nil {
			return err
		}

		err = s.db.RecordReleases(providers.GithubReleaseProviderName, sub.Group, sub.Repo, releases)
		if err != nil {
			return err
		}
	}

	return nil
}