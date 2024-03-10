package main

import (
	"github.com/nerzhul/releasetracker/pkg/releasetracker"
	"github.com/nerzhul/releasetracker/pkg/utils"
)

func main() {
	httpPort := utils.GetIntEnvOrDefault("HTTP_PORT", 8080)
	tracker := releasetracker.NewReleaseTracker(httpPort)
	err := tracker.Run()
	if err != nil {
		panic(err)
	}
}
