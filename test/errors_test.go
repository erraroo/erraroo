package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/erraroo/erraroo/api/events"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	account, err := models.Accounts.Create()
	assert.Nil(t, err)
	assert.False(t, 0 == account.ID, "should not be 0")
}

func TestCreateEvent(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)

	request := events.CreateEventRequest{
		Kind: "js.error",
		Data: map[string]interface{}{
			"message": "error thrown",
		},
	}
	req, res := rr("POST", "/api/v1/events", request)
	req.Header.Set("X-Token", project.Token)
	_app.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)

	errors, err := models.Events.ListForProject(project)
	assert.Nil(t, err)
	assert.NotEmpty(t, errors)

	e := errors[0]
	assert.NotEmpty(t, e.Checksum, "the checksum was generated")

	jobs.Work(_app)

	// It should have created a group for the project
	groups, err := models.Groups.FindQuery(models.GroupQuery{
		ProjectID: project.ID,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, groups.Groups)
	assert.Equal(t, groups.Groups[0].Checksum, e.Checksum)
}

func TestEventShow(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)
	e, _ := models.Events.Create(project.Token, "{}")

	req, res := rr("GET", fmt.Sprintf("/api/v1/errors/%d", e.ID), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.ShowEvent{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, e.ID, response.Event.ID)
	assert.Equal(t, e.Payload, response.Event.Payload)
	assert.Equal(t, e.Checksum, response.Event.Checksum)
}

func TestEventShowOnlyShowsEventsOwnedByUser(t *testing.T) {
	account2, _ := models.Accounts.Create()
	project, _ := models.Projects.Create("test project", account2.ID)
	e, _ := models.Events.Create(project.Token, "{}")

	req, res := rr("GET", fmt.Sprintf("/api/v1/errors/%d", e.ID), nil)
	req.Header.Add("Authorization", _token)
	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestEventsByProjectId(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)
	e, _ := models.Events.Create(project.Token, "{}")
	group, _ := models.Groups.FindOrCreate(project, e)

	req, res := rr("GET", fmt.Sprintf("/api/v1/errors?project_id=%d&group_id=%d", project.ID, group.ID), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.Events{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, len(response.Events), 1)
	assert.Equal(t, e.Payload, response.Events[0].Payload)
	assert.Equal(t, e.Checksum, response.Events[0].Checksum)
}
