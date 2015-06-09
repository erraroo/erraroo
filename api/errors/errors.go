package errors

import (
	"net/http"
	"strconv"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

func Show(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	id, err := cx.GetID(r)
	if err != nil {
		return err
	}

	e, err := models.Errors.FindByID(id)
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

	return cx.JSON(w, http.StatusOK, serializers.NewShowError(e))
}

func Index(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	projectID, err := cx.StrToID(r.URL.Query().Get("project_id"))
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

	groupID, err := cx.StrToID(r.URL.Query().Get("group_id"))
	if err != nil {
		return err
	}

	group, err := models.Groups.FindByID(groupID)
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 50
	}

	errs, err := models.Errors.FindQuery(models.ErrorQuery{
		ProjectID: project.ID,
		Checksum:  group.Checksum,
		QueryOptions: models.QueryOptions{
			Page:    page,
			PerPage: limit,
		},
	})
	if err != nil {
		return err
	}

	return cx.JSON(w, http.StatusOK, serializers.NewErrors(errs))
}
