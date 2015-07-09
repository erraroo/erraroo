package api

import (
	"net/http"

	"github.com/erraroo/erraroo/api/bus"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

// UpdateErrorRequest incoming update request
type UpdateErrorRequest struct {
	Error ErrorParams
}

// ErrorParams the params that we can safely assign to a Error
type ErrorParams struct {
	Resolved bool
	Muted    bool
}

// ErrorsIndex returns the paginated errors filtered by a project_id
func ErrorsIndex(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	projectID, err := QueryToID(r, "project_id")
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

	query := models.ErrorQuery{}
	query.QueryOptions = QueryOptions(r)
	query.PerPage = 50
	query.ProjectID = project.ID
	query.Status = r.URL.Query().Get("status")

	if id, err := StrToID(r.URL.Query().Get("library")); err == nil {
		query.Libaries = []int64{id}
	}

	groups, err := models.Errors.FindQuery(query)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewErrors(groups))
}

// ErrorsUpdate updates the error record with an incoming UpdateErrorRequest
func ErrorsUpdate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	group, err := getAuthorizedError(r, ctx)
	if err != nil {
		return err
	}

	request := UpdateErrorRequest{}
	Decode(r, &request)

	group.Muted = request.Error.Muted
	group.Resolved = request.Error.Resolved
	err = models.Errors.Update(group)
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(group.ProjectID)
	if err != nil {
		return err
	}

	payload := serializers.NewUpdateError(project, group)

	err = bus.Push(project.Channel(), bus.Notifcation{
		Name:    "errors.update",
		Payload: payload,
	})

	if err != nil {
		logger.Error("bus.push", "err", err)
	}

	return JSON(w, http.StatusOK, payload)
}

// ErrorsShow returns the full group
func ErrorsShow(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	group, err := getAuthorizedError(r, ctx)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowError(group))
}

func getAuthorizedError(r *http.Request, ctx *cx.Context) (*models.Error, error) {
	id, err := GetID(r)
	if err != nil {
		return nil, err
	}

	group, err := models.Errors.FindByID(id)
	if err != nil {
		return nil, err
	}

	project, err := models.Projects.FindByID(group.ProjectID)
	if err != nil {
		return nil, err
	}

	if !ctx.User.CanSee(project) {
		return nil, models.ErrNotFound
	}

	return group, nil
}
