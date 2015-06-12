package timings

import (
	"net/http"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

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

	timings, err := models.Timings.Last7Days(project)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewTimings(timings))
}
