package releasesyncer

import (
	"context"

	"github.com/nerzhul/releasetracker/pkg/providers"
	"github.com/nerzhul/releasetracker/pkg/release"
)

type ReleaseSyncerDB struct {
	*providers.DatabaseReleaseProvider
}

func NewReleaseSyncerDB(databaseURL string) *ReleaseSyncerDB {
	return &ReleaseSyncerDB{
		providers.NewDatabaseReleaseProvider(databaseURL),
	}
}

func (db *ReleaseSyncerDB) getSubscribedReleases(providerName string) ([]release.ReleaseMetadata, error) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT group_name, repo_name FROM releases_subscriptions
		WHERE provider_name = $1`, providerName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var releases []release.ReleaseMetadata

	for rows.Next() {
		release := release.ReleaseMetadata{
			Provider: providerName,
		}

		err = rows.Scan(&release.Group, &release.Repo)
		if err != nil {
			return nil, err
		}

		releases = append(releases, release)
	}

	return releases, nil
}