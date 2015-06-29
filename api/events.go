package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/erraroo/erraroo/usecases"
)

const rateLimitDuration = 60 * time.Second

type CreateEventRequest struct {
	Kind string                 `json:"kind"`
	Data map[string]interface{} `json:"data"`
}

func EventsCreate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	token := r.Header.Get("X-Token")
	if token == "" {
		return errors.New("token was blank")
	}

	plan, err := models.Plans.FindByToken(token)
	if err != nil {
		return err
	}

	ok, err := Limiter.Check(token, rateLimitDuration, plan.RequestsPerMinute)
	if err != nil {
		logger.Error("Limiter.Check", "err", err)
		return err
	}

	if !ok {
		w.WriteHeader(420)
		logger.Error("rate limit exceeded", "token", token)
		return usecases.RateExceeded(token)
	}

	request := CreateEventRequest{}
	Decode(r, &request)

	payload, err := json.Marshal(request.Data)
	if err != nil {
		return err
	}
	data := string(payload)

	event, err := models.Events.Create(token, request.Kind, data)
	if err == models.ErrNotFound {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	if err != nil {
		return err
	}

	if event.IsAsync() {
		payload, err := json.Marshal(event.ID)
		if err != nil {
			return err
		}

		err = jobs.Push("event.process", payload)
		if err != nil {
			return err
		}
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func EventsShow(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	id, err := GetID(r)
	if err != nil {
		return err
	}

	e, err := models.Events.FindByID(id)
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(e.ProjectID)
	if err != nil {
		return err
	}

	if !ctx.User.CanSee(project) {
		return models.ErrNotFound
	}

	return JSON(w, http.StatusOK, serializers.NewShowEvent(e))
}

func EventsIndex(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	projectID, err := QueryToID(r, "project_id")
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

	events, err := models.Events.FindQuery(models.EventQuery{
		Checksum:     r.URL.Query().Get("checksum"),
		Kind:         r.URL.Query().Get("kind"),
		ProjectID:    project.ID,
		QueryOptions: QueryOptions(r),
	})

	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewEvents(events))
}
