package main

import (
	"github.com/nerzhul/releasetracker/pkg/releasetracker"
	"github.com/nerzhul/releasetracker/pkg/utils"
)

func main() {
	httpPort := utils.GetIntEnvOrDefault("HTTP_PORT", 8080)
	databaseURL := utils.GetEnvOrDefault("DATABASE_URL", "postgres://postgres:password@localhost:5432/releasetracker?sslmode=disable")

	tracker := releasetracker.NewReleaseTracker(httpPort, databaseURL)
	err := tracker.Run()
	if err != nil {
		panic(err)
	}
}
