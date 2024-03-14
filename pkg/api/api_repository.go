package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *ReleaseTrackerAPI) PostV1RepoSubscribeProviderGroupRepo(c echo.Context, provider string, group string, repo string) error {
	log := a.log.WithValues("provider", provider, "group", group, "repo", repo)
	log.Info("Subscribing to repository")

	sub, err := a.db.HasSubscribedToReleases(provider, group, repo)
	if err != nil {
		log.Error(err, "Failed to check if already subscribed to repository")
		return c.JSON(http.StatusInternalServerError, StatusOnlyReponse{
			Status: StatusResponseServerError,
		})
	}

	if sub {
		log.Info("Already subscribed to repository")
		return c.JSON(http.StatusConflict, StatusOnlyReponse{
			Status: StatusResponseConflict,
		})
	}

	if err := a.db.SubscribeReleases(provider, group, repo); err != nil {
		log.Error(err, "Failed to subscribe to repository")
		return c.JSON(http.StatusInternalServerError, StatusOnlyReponse{
			Status: StatusResponseServerError,
		})
	}
	return c.JSON(http.StatusCreated, StatusOnlyReponse{
		Status: StatusResponseOK,
	})
}
