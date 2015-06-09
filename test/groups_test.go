package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/erraroo/erraroo/api/groups"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/stretchr/testify/assert"
)

func TestQueryGroups(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)
	e, _ := models.Errors.Create(project.Token, "{}")
	group, err := models.Groups.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group)

	req, res := rr("GET", fmt.Sprintf("/api/v1/groups?project_id=%d", project.ID), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.Groups{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, len(response.Groups), 1)
	assert.Equal(t, response.Meta.Pagination.Page, 1)
	assert.Equal(t, response.Meta.Pagination.Pages, 1)
	assert.Equal(t, response.Meta.Pagination.Limit, 50)
	assert.Equal(t, response.Meta.Pagination.Total, 1)
	assert.Equal(t, response.Groups[0].ID, group.ID)
	assert.Equal(t, response.Groups[0].Checksum, group.Checksum)
	assert.Equal(t, response.Groups[0].Message, group.Message)
	assert.Equal(t, response.Groups[0].ProjectID, group.ProjectID)

	req, res = rr("GET", fmt.Sprintf("/api/v1/groups?project_id=%d", 0), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, res.Code, http.StatusNotFound)
}

func TestUpdateGroups(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)
	e, _ := models.Errors.Create(project.Token, "{}")
	group, err := models.Groups.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group)

	request := groups.UpdateGroupRequest{}
	request.Group.Resolved = true

	req, res := rr("PUT", fmt.Sprintf("/api/v1/groups/%d", group.ID), request)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.ShowGroup{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, response.Group.Resolved, true)
}
