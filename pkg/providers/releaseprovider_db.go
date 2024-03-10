package providers

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerzhul/releasetracker/pkg/release"
)

type DatabaseReleaseProvider struct {
	pool *pgxpool.Pool
}

func NewDatabaseReleaseProvider(databaseURL string) *DatabaseReleaseProvider {
	// TODO: think about db migration
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	c, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to acquire connection from pool: %v\n", err)
		os.Exit(1)
	}

	if err := c.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	return &DatabaseReleaseProvider{
		pool: pool,
	}
}

func (p *DatabaseReleaseProvider) GetReleases(group string, repo string, maxReleases int) (*release.ReleaseList, error) {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	rows, err := c.Query(context.Background(), "SELECT release_tag, release_time FROM releases WHERE group = $1 AND repo = $2 ORDER BY release_time DESC LIMIT $3",
		group, repo, maxReleases)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	releaseList := &release.ReleaseList{
		Group:    group,
		Repo:     repo,
		Releases: make([]release.Release, 0),
	}

	for rows.Next() {
		release := release.Release{}

		err = rows.Scan(&release.Version, &release.ReleaseDate)
		if err != nil {
			return nil, err
		}

		releaseList.Releases = append(releaseList.Releases, release)
	}

	return releaseList, nil
}

func (p *DatabaseReleaseProvider) RecordReleases(group string, repo string, releases *release.ReleaseList) error {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	for _, release := range releases.Releases {
		_, err := c.Exec(context.Background(), "INSERT INTO releases (group, repo, release_tag, release_time) VALUES ($1, $2, $3, $4)" +
			"ON CONFLICT (group, repo, release_tag) DO UPDATE SET release_time = $4",
			group, repo, release.Version, release.ReleaseDate)
		if err != nil {
			return err
		}
	}

	return nil
}
