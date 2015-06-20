package errors

import (
	"net/http"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

func Show(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	id, err := api.GetID(r)
	if err != nil {
		return err
	}

	e, err := models.Events.FindByID(id)
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(e.ProjectID)
	if err != nil {
		return err
	}

	if !ctx.User.CanSee(project) {
		return models.ErrNotFound
	}

	return api.JSON(w, http.StatusOK, serializers.NewShowEvent(e))
}

func Index(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	projectID, err := api.QueryToID(r, "project_id")
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(projectID)
	if err != nil {
		return err
	}

	if !ctx.User.CanSee(project) {
		return models.ErrNotFound
	}

	events, err := models.Events.FindQuery(models.EventQuery{
		Checksum:     r.URL.Query().Get("checksum"),
		ProjectID:    project.ID,
		QueryOptions: api.QueryOptions(r),
	})

	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewEvents(events))
}
