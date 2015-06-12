package projects

import (
	"net/http"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

type ProjectRequest struct {
	Project ProjectParams
}

type ProjectParams struct {
	Name string
}

func Create(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	params := projectParams(r)
	project, err := models.Projects.Create(params.Name, ctx.User.AccountID)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusCreated, serializers.NewShowProject(project))
}

func Index(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	projects, err := models.Projects.ByAccountID(ctx.User.AccountID)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewProjects(projects))
}

func Show(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewShowProject(project))
}

// Update updates the project record with an incoming UpdateProjectRequest
func Update(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	params := projectParams(r)
	project.Name = params.Name
	err = models.Projects.Update(project)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, serializers.NewShowProject(project))
}

func getAuthorizedProject(r *http.Request, ctx *cx.Context) (*models.Project, error) {
	id, err := api.GetID(r)
	if err != nil {
		return nil, err
	}

	project, err := models.Projects.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !ctx.User.CanSee(project) {
		return nil, models.ErrNotFound
	}

	return project, nil
}

func projectParams(r *http.Request) ProjectParams {
	request := ProjectRequest{}
	api.Decode(r, &request)
	return request.Project
}
