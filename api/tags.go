package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
)

func ErrorTagsIndex(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	errorID, err := GetID(r)
	if err != nil {
		return err
	}

	e, err := models.Errors.FindByID(errorID)
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

	return JSON(w, http.StatusOK, struct {
		Tags []models.TagValue
	}{Tags: e.Tags})
}

// ErrorsUpdate updates the error record with an incoming UpdateErrorRequest
