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

func NewDatabaseReleaseProvider() *DatabaseReleaseProvider {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	c, err := pool.Acquire(context.Background())
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

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		return nil, err
	}

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
