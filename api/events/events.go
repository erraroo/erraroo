package events

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/usecases"
)

const rateLimitDuration = 60 * time.Second

type CreateEventRequest struct {
	Kind string                 `json:"kind"`
	Data map[string]interface{} `json:"data"`
}

func Create(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	token := r.Header.Get("X-Token")
	if token == "" {
		return errors.New("token was blank")
	}

	plan, err := models.Plans.FindByToken(token)
	if err != nil {
		return err
	}

	ok, err := api.Limiter.Check(token, rateLimitDuration, plan.RequestsPerMinute)
	if err != nil {
		logger.Error("api.Limiter.Check", "err", err)
		return err
	}

	if !ok {
		w.WriteHeader(420)
		logger.Error("rate limit exceeded", "token", token)
		return usecases.RateExceeded(token)
	}

	request := CreateEventRequest{}
	api.Decode(r, &request)

	payload, err := json.Marshal(request.Data)
	if err != nil {
		return err
	}
	data := string(payload)

	switch request.Kind {
	case "js.error":
		e, err := models.Events.Create(token, data)
		if err == models.ErrNotFound {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		if err != nil {
			return err
		}

		payload, err := json.Marshal(e.ID)
		if err != nil {
			return err
		}

		err = jobs.Push("create.error", payload)
		if err != nil {
			return err
		}
	case "js.timing":
		_, err := models.Timings.Create(token, data)
		if err == models.ErrNotFound {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		if err != nil {
			return err
		}
	case "js.log":
		logger.Info("js.log", "payload", data)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
