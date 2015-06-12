package groups

import (
	"net/http"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

// UpdateGroupRequest incoming update request
type UpdateGroupRequest struct {
	Group GroupParams
}

// GroupParams the params that we can safely assign to a Group
type GroupParams struct {
	Resolved bool
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

	query := models.GroupQuery{}
	query.PerPage = 50
	query.ProjectID = project.ID
	query.QueryOptions.Page = api.Page(r)

	groups, err := models.Groups.FindQuery(query)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewGroups(groups))
}

// Update updates the group record with an incoming UpdateGroupRequest
func Update(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	group, err := getAuthorizedGroup(r, ctx)
	if err != nil {
		return err
	}

	request := UpdateGroupRequest{}
	api.Decode(r, &request)

	group.Resolved = request.Group.Resolved
	err = models.Groups.Update(group)
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(group.ProjectID)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewUpdateGroup(project, group))
}

// Show returns the full group
func Show(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	group, err := getAuthorizedGroup(r, ctx)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewShowGroup(group))
}

func getAuthorizedGroup(r *http.Request, ctx *cx.Context) (*models.Group, error) {
	id, err := api.GetID(r)
	if err != nil {
		return nil, err
	}

	group, err := models.Groups.FindByID(id)
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
