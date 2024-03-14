package api

import (
	"github.com/go-logr/logr"
	"github.com/nerzhul/releasetracker/pkg/providers"
)

type ReleaseTrackerAPI struct {
	log logr.Logger
	db  *providers.DatabaseReleaseProvider
}

func NewReleaseTrackerAPI(logger logr.Logger, databaseURL string) *ReleaseTrackerAPI {
	return &ReleaseTrackerAPI{
		log: logger,
		db:  providers.NewDatabaseReleaseProvider(databaseURL),
	}
}
