package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/models/events"
	"github.com/erraroo/erraroo/serializers"
	"github.com/erraroo/erraroo/usecases"
)

const (
	rateLimitDuration             = 60 * time.Second
	notificationRateLimitDuration = 30 * time.Minute
	notificationRateLimitMax      = 1
)

const (
	StatusSlowYourRoll = 420
)

var (
	planCache = models.NewPlanCache()
)

func EventsCreate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	token := r.Header.Get("X-Token")
	if token == "" {
		return errors.New("token was blank")
	}

	plan, err := planCache.FindByToken(token)
	if err != nil {
		return err
	}

	ok, err := Limiter.Check(token, rateLimitDuration, plan.RequestsPerMinute)
	if err != nil {
		logger.Error("Limiter.Check", "key", token, "err", err)
		return err
	}

	if !ok {
		w.WriteHeader(StatusSlowYourRoll)
		return tryNotify(token)
	}

	request := events.CreateEventRequest{}
	err = Decode(r, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	err = events.Ingest(token, request)
	if err == models.ErrNotFound {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	if err != nil {
		return err
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

func tryNotify(token string) error {
	key := token + "-exceeded"

	ok, err := Limiter.Check(key, 10*time.Second, 1)
	if err != nil {
		logger.Error("checking notifcation rate limit", "key", key, "err", err)
		return err
	}

	if !ok {
		logger.Debug("already notified", "key", key)
		return nil
	} else {
		return usecases.RateExceeded(token)
	}
}
