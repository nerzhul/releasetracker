package main

import (
	"os"

	"github.com/nerzhul/releasetracker/pkg/releasesyncer"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL environment variable is not set")
	}

	syncer := releasesyncer.NewReleaseSyncer(databaseURL)
	err := syncer.SyncReleases()
	if err != nil {
		panic(err)
	}
}
