package api

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/erraroo/erraroo/config"
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

func (p ProjectParams) Validate() (models.ValidationErrors, error) {
	errs := models.NewValidationErrors()
	if p.Name == "" {
		errs.Add("Name", "can't be blank")
	}

	return errs, nil
}

func ProjectsCreate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	params := projectParams(r)
	errors, err := params.Validate()
	if err != nil {
		return err
	}

	if errors.Any() {
		return errors
	}

	project, err := models.Projects.Create(params.Name, ctx.User.AccountID)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusCreated, serializers.NewShowProject(project))
}

func ProjectsIndex(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	projects, err := models.Projects.ByAccountID(ctx.User.AccountID)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewProjects(projects))
}

func ProjectsShow(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowProject(project))
}

// ProjectsUpdate updates the project record with an incoming UpdateProjectRequest
func ProjectsUpdate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	params := projectParams(r)
	errors, err := params.Validate()
	if err != nil {
		return err
	}

	if errors.Any() {
		return errors
	}

	project.Name = params.Name
	err = models.Projects.Update(project)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowProject(project))
}

func ProjectsDelete(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	err = models.Projects.Delete(project)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, nil)
}

func ProjectsRepository(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	repository, err := models.FindRepositoryByProjectID(project.ID)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowRepository(repository))
}

func ProjectsOutdatedRevisions(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	revisions, err := models.FindOutdatedRevisionsByProjectID(project.ID)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, map[string]interface{}{
		"OutdatedRevisions": revisions,
	})
}

func ProjectsRegenerateToken(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	project, err := getAuthorizedProject(r, ctx)
	if err != nil {
		return err
	}

	token, err := models.Projects.GenerateToken()
	if err != nil {
		return err
	}

	project.Token = token

	err = models.Projects.Update(project)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowProject(project))
}

type SignedProjectTokenRequest struct {
	SignedProjectToken struct {
		ProjectID string
		Token     string
	}
}

func SignedProjectTokens(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := SignedProjectTokenRequest{}
	Decode(r, &request)

	log.Println(request)

	id, err := StrToID(request.SignedProjectToken.ProjectID)
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(id)
	if err != nil {
		return err
	}

	if !ctx.User.CanSee(project) {
		return models.ErrNotFound
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["project_id"] = project.ID
	tokenString, err := token.SignedString(config.TokenSigningKey)
	if err != nil {
		return err
	}

	request.SignedProjectToken.Token = tokenString

	return JSON(w, http.StatusCreated, request)
}

func getAuthorizedProject(r *http.Request, ctx *cx.Context) (*models.Project, error) {
	id, err := GetID(r)
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
	Decode(r, &request)
	return request.Project
}
