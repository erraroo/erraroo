package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/erraroo/erraroo/usecases"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	account, err := models.CreateAccount()
	assert.Nil(t, err)
	assert.False(t, 0 == account.ID, "should not be 0")
}

func TestCreateEvent(t *testing.T) {
	project, err := models.Projects.Create("test project", _account.ID)
	assert.Nil(t, err)

	plan, err := models.Plans.Create(_account, "default")
	assert.Nil(t, err)
	assert.NotNil(t, plan)

	request := usecases.CollectEventRequest{
		Kind: "js.error",
		Data: map[string]interface{}{
			"message": "error thrown",
		},
	}
	req, res := rr("POST", "/api/v1/events", request)
	req.Header.Set("X-Token", project.Token)
	_app.ServeHTTP(res, req)
	jobs.Work(_app)

	assert.Equal(t, http.StatusCreated, res.Code)

	events, err := models.Events.FindQuery(models.EventQuery{ProjectID: project.ID})
	assert.Nil(t, err)
	assert.NotEmpty(t, events.Events)

	e := events.Events[0]
	assert.NotEmpty(t, e.Checksum, "the checksum was generated")

	jobs.Work(_app)

	// It should have created a group for the project
	groups, err := models.Errors.FindQuery(models.ErrorQuery{
		ProjectID: project.ID,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, groups.Errors)
	assert.Equal(t, groups.Errors[0].Checksum, e.Checksum)
}

func TestEventShow(t *testing.T) {
	project, _ := models.Projects.Create("test project", _account.ID)
	event := models.NewEvent(project, "js.error", "{}")
	err := models.Events.Insert(event)
	assert.Nil(t, err)

	req, res := rr("GET", fmt.Sprintf("/api/v1/events/%d", event.ID), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.ShowEvent{}
	json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, event.ID, response.Event.ID)
	//assert.Equal(t, event.Payload, response.Event.Payload)
	assert.Equal(t, event.Checksum, response.Event.Checksum)
}

func TestEventShowOnlyShowsEventsOwnedByUser(t *testing.T) {
	account2, _ := models.CreateAccount()
	project, _ := models.Projects.Create("test project", account2.ID)
	event := models.NewEvent(project, "js.error", "{}")
	err := models.Events.Insert(event)
	assert.Nil(t, err)

	req, res := rr("GET", fmt.Sprintf("/api/v1/events/%d", event.ID), nil)
	req.Header.Add("Authorization", _token)
	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestEventsByProjectId(t *testing.T) {
	project, err := models.Projects.Create("test project", _account.ID)
	assert.Nil(t, err)

	e := models.NewEvent(project, "js.error", "{}")
	err = models.Events.Insert(e)
	assert.Nil(t, err)

	group, err := models.Errors.FindOrCreate(project, e)
	assert.Nil(t, err)

	req, res := rr("GET", fmt.Sprintf("/api/v1/events?project_id=%d&group_id=%d", project.ID, group.ID), nil)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	response := serializers.Events{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, len(response.Events), 1)
	//assert.Equal(t, e.Payload, response.Events[0].Payload)
	assert.Equal(t, e.Checksum, response.Events[0].Checksum)
}

type alwaysLimter struct{}

func (a alwaysLimter) Check(key string, d time.Duration, count int) (bool, error) {
	return false, nil
}

func TestCreateEventIsRateLimited(t *testing.T) {
	api.Limiter = alwaysLimter{}
	defer func() {
		api.Limiter = api.NoLimiter()
	}()

	project, err := models.Projects.Create("test project", _account.ID)
	assert.Nil(t, err)

	req, res := rr("POST", "/api/v1/events", nil)
	req.Header.Set("X-Token", project.Token)
	_app.ServeHTTP(res, req)
	assert.Equal(t, api.StatusSlowYourRoll, res.Code)
}

type onceLimter struct{ count int }

func (o *onceLimter) Check(key string, d time.Duration, count int) (bool, error) {
	if o.count == 0 {
		log.Printf("limiting %s\n", key)
		o.count++
		return false, nil
	}

	log.Printf("allowing %s\n", key)
	return true, nil
}

func TestCreateEventIsRateLimitedNotifcationFiresOnce(t *testing.T) {
	api.Limiter = &onceLimter{}
	defer func() {
		api.Limiter = api.NoLimiter()
	}()

	project, err := models.Projects.Create("test project", _account.ID)
	assert.Nil(t, err)

	req, res := rr("POST", "/api/v1/events", nil)
	req.Header.Set("X-Token", project.Token)
	_app.ServeHTTP(res, req)

	notified, err := models.RateLimitNotifcations.WasRecentlyNotified(_account)
	assert.Nil(t, err)
	assert.True(t, notified)
}
