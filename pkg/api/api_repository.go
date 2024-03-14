package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *ReleaseTrackerAPI) PostV1RepoSubscribeProviderGroupRepo(c echo.Context, provider string, group string, repo string) error {
	log := a.log.WithValues("provider", provider, "group", group, "repo", repo)
	log.Info("Subscribing to repository")
	
	resp := StatusOnlyReponse{
		Status: "OK",
	}

	return c.JSON(http.StatusCreated, resp)
}
