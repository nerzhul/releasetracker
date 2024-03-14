package api

import "github.com/go-logr/logr"

type ReleaseTrackerAPI struct {
	log logr.Logger
}

func NewReleaseTrackerAPI(logger logr.Logger) *ReleaseTrackerAPI {
	return &ReleaseTrackerAPI{
		log: logger,
	}
}
