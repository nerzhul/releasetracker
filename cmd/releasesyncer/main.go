package main

import (
	"github.com/nerzhul/releasetracker/pkg/releasesyncer"
	"github.com/nerzhul/releasetracker/pkg/utils"
)

func main() {
	databaseURL := utils.GetEnvOrDefault("DATABASE_URL", "postgres://postgres:password@localhost:5432/releasetracker?sslmode=disable")

	syncer := releasesyncer.NewReleaseSyncer(databaseURL)
	err := syncer.SyncReleases()
	if err != nil {
		panic(err)
	}
}
