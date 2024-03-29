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

func (p *DatabaseReleaseProvider) Conn(ctx context.Context) (*pgxpool.Conn, error) {
	return p.pool.Acquire(ctx)
}

func (p *DatabaseReleaseProvider) GetReleases(providerName string, group string, repo string, maxReleases int) (*release.ReleaseList, error) {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	defer c.Release()

	providerID, err := p.getProviderID(providerName, group, repo)
	if err != nil {
		return nil, err
	}

	releaseList := &release.ReleaseList{
		Metadata: release.ReleaseMetadata{
			Provider: providerName,
			Group:    group,
			Repo:     repo,
		},
		Releases: make([]release.Release, 0),
	}

	rows, err := c.Query(context.Background(), "SELECT release_version, release_time FROM releases WHERE provider_id = $1 ORDER BY release_time DESC LIMIT $2",
		providerID, maxReleases)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func (p *DatabaseReleaseProvider) SubscribeReleases(provider string, group string, repo string) error {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer c.Release()

	providerID, err := p.getOrCreateProviderID(provider, group, repo)
	if err != nil {
		return err
	}

	_, err = c.Exec(context.Background(), "INSERT INTO releases_subscriptions (provider_id) VALUES ($1) ON CONFLICT DO NOTHING",
		providerID)
	if err != nil {
		return err
	}

	return nil
}

func (p *DatabaseReleaseProvider) HasSubscribedToReleases(provider string, group string, repo string) (bool, error) {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return false, err
	}

	defer c.Release()

	providerID, err := p.getProviderID(provider, group, repo)
	if err != nil {
		return false, err
	}

	if providerID == "" {
		return false, nil
	}

	rows, err := c.Query(context.Background(), "SELECT provider_id FROM releases_subscriptions WHERE provider_id = $1",
		providerID)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), nil
}

func (p *DatabaseReleaseProvider) RecordReleases(provider string, group string, repo string, releases *release.ReleaseList) error {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer c.Release()

	providerID, err := p.getOrCreateProviderID(provider, group, repo)
	if err != nil {
		return err
	}

	for _, release := range releases.Releases {
		_, err := c.Exec(context.Background(), "INSERT INTO releases (provider_id, release_tag, release_time) VALUES ($1, $2, $3)"+
			"ON CONFLICT (provider_id, release_tag) DO UPDATE SET release_time = $3",
			providerID, release.Version, release.ReleaseDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DatabaseReleaseProvider) getProviderID(provider string, group string, repo string) (string, error) {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return "", err
	}

	defer c.Release()

	provRows, err := c.Query(context.Background(), "SELECT id FROM release_providers WHERE provider = $1 AND group_name = $2 AND repo_name = $3",
		provider, group, repo)

	if err != nil {
		return "", err
	}

	defer provRows.Close()

	if !provRows.Next() {
		return "", nil
	}

	var providerID string

	err = provRows.Scan(&providerID)
	if err != nil {
		return "", err
	}

	return providerID, nil
}

func (p *DatabaseReleaseProvider) createProvider(provider string, group string, repo string) (string, error) {
	c, err := p.pool.Acquire(context.Background())
	if err != nil {
		return "", err
	}

	rows, err := c.Query(context.Background(), "INSERT INTO release_providers (provider_name, group_name, repo_name) " +
		"VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id",
		provider, group, repo)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	if rows.Next() {
		var providerID string
		err = rows.Scan(&providerID)
		if err != nil {
			return "", err
		}

		return providerID, nil
	}

	return "", fmt.Errorf("failed to create release provider")
}

func (p *DatabaseReleaseProvider) getOrCreateProviderID(provider string, group string, repo string) (string, error) {
	providerID, err := p.getProviderID(provider, group, repo)
	if err != nil {
		return "", err
	}

	if providerID == "" {
		return p.createProvider(provider, group, repo)
	}

	return providerID, nil
}
