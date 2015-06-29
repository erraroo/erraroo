package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	request := api.ProjectRequest{}
	request.Project.Name = "test project"

	req, res := rr("POST", "/api/v1/projects", request)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)

	response := serializers.ShowProject{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "test project", response.Project.Name)
	assert.Equal(t, _user.AccountID, response.Project.AccountID)
	assert.NotEmpty(t, response.Project.Token)
}

func TestProjectIndex(t *testing.T) {
	models.Projects.Create("test project", _account.ID)

	req, res := rr("GET", "/api/v1/projects", nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.Projects{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotEmpty(t, response.Projects)
}

func TestProjectShow(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)

	req, res := rr("GET", fmt.Sprintf("/api/v1/projects/%d", project.ID), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.ShowProject{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, project.ID, response.Project.ID)
	assert.Equal(t, project.Name, response.Project.Name)
	assert.Equal(t, project.Token, response.Project.Token)
}
