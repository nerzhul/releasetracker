package releasetracker

import (
	"fmt"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nerzhul/releasetracker/pkg/utils"
)

type ReleaseTracker struct {
	listerningPort int
}

func NewReleaseTracker(port int) *ReleaseTracker {
	return &ReleaseTracker{
		listerningPort: port,
	}
}

func (r *ReleaseTracker) Run() error {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(utils.CreateHTTPLoggerConfig()))
	e.Use(middleware.Recover())

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	// TODO handlers
	return e.Start(fmt.Sprintf(":%d", r.listerningPort))
}
