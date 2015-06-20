package groups

import (
	"net/http"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/cx"
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

// Index returns the paginated groups filtered by a project_id
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

	query := models.ErrorQuery{}
	query.PerPage = 50
	query.ProjectID = project.ID
	query.QueryOptions.Page = api.Page(r)

	groups, err := models.Errors.FindQuery(query)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewErrors(groups))
}

// Update updates the group record with an incoming UpdateErrorRequest
func Update(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	group, err := getAuthorizedError(r, ctx)
	if err != nil {
		return err
	}

	request := UpdateErrorRequest{}
	api.Decode(r, &request)

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

	return api.JSON(w, http.StatusOK, serializers.NewUpdateError(project, group))
}

// Show returns the full group
func Show(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	group, err := getAuthorizedError(r, ctx)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewShowError(group))
}

func getAuthorizedError(r *http.Request, ctx *cx.Context) (*models.Error, error) {
	id, err := api.GetID(r)
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
