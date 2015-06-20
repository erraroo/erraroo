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

func TestQueryErrors(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)
	e, _ := models.Events.Create(project.Token, "{}")
	group, err := models.Errors.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group)

	req, res := rr("GET", fmt.Sprintf("/api/v1/errors?project_id=%d", project.ID), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.Errors{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Meta.Pagination.Page, 1)
	assert.Equal(t, response.Meta.Pagination.Pages, 1)
	assert.Equal(t, response.Meta.Pagination.Limit, 50)
	assert.Equal(t, response.Meta.Pagination.Total, 1)
	assert.Equal(t, response.Errors[0].ID, group.ID)
	assert.Equal(t, response.Errors[0].Checksum, group.Checksum)
	assert.Equal(t, response.Errors[0].Message, group.Message)
	assert.Equal(t, response.Errors[0].ProjectID, group.ProjectID)

	req, res = rr("GET", fmt.Sprintf("/api/v1/errors?project_id=%d", 0), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, res.Code, http.StatusNotFound)
}

func TestUpdateErrors(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)
	e, err := models.Events.Create(project.Token, "{}")
	assert.Nil(t, err)

	group, err := models.Errors.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group)

	request := api.UpdateErrorRequest{}
	request.Error.Resolved = true

	req, res := rr("PUT", fmt.Sprintf("/api/v1/errors/%d", group.ID), request)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.ShowError{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, response.Error.Resolved, true)
}
