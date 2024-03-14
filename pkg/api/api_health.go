package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *ReleaseTrackerAPI) GetV1HealthReadyz(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
