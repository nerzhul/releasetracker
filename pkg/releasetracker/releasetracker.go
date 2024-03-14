package releasetracker

import (
	"fmt"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nerzhul/releasetracker/pkg/api"
	"github.com/nerzhul/releasetracker/pkg/utils"

	"github.com/go-logr/zapr"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
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
	zapLog, err := zap.NewProduction()
	if err != nil {
		return err
	}

	logger := zapr.NewLogger(zapLog)
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(utils.CreateHTTPLoggerConfig()))
	e.Use(middleware.Recover())

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	myAPI := api.NewReleaseTrackerAPI(logger)
	api.RegisterHandlers(e, myAPI)

	// TODO: CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	return e.Start(fmt.Sprintf(":%d", r.listerningPort))
}
