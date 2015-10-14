package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

type RepositoryRequest struct {
	Repository RepositoryParams
}

type RepositoryParams struct {
	GithubOrg  string
	GithubRepo string
}

func RepositoriesUpdate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	repository, err := getAuthorizedRepository(r, ctx)
	if err != nil {
		return err
	}

	params := repositoryParams(r)
	repository.GithubOrg = params.GithubOrg
	repository.GithubRepo = params.GithubRepo

	err = models.InsertRepository(repository)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowRepository(repository))
}

func RepositoriesDelete(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	repository, err := getAuthorizedRepository(r, ctx)
	if err != nil {
		return err
	}

	params := repositoryParams(r)
	repository.GithubOrg = params.GithubOrg
	repository.GithubRepo = params.GithubRepo

	err = models.DeleteRepository(repository)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func repositoryParams(r *http.Request) RepositoryParams {
	request := RepositoryRequest{}
	Decode(r, &request)
	return request.Repository
}

func getAuthorizedRepository(r *http.Request, ctx *cx.Context) (*models.Repository, error) {
	id, err := GetID(r)
	if err != nil {
		return nil, err
	}

	repository, err := models.FindRepositoryByID(id)
	if err != nil {
		return nil, err
	}

	project, err := models.Projects.FindByID(repository.ProjectID)
	if err != nil {
		return nil, err
	}

	if !ctx.User.CanSee(project) {
		return nil, models.ErrNotFound
	}

	return repository, nil
}
