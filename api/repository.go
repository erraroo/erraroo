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
	id, err := GetID(r)
	if err != nil {
		return err
	}

	repository, err := models.FindRepositoryByID(id)
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(repository.ProjectID)
	if err != nil {
		return err
	}

	if !ctx.User.CanSee(project) {
		return models.ErrNotFound
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

func repositoryParams(r *http.Request) RepositoryParams {
	request := RepositoryRequest{}
	Decode(r, &request)
	return request.Repository
}
