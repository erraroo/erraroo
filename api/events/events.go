package events

import (
	"encoding/json"
	"net/http"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
)

type CreateEventRequest struct {
	Kind  string                 `json:"kind"`
	Token string                 `json:"token"`
	Data  map[string]interface{} `json:"data"`
}

func Create(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := CreateEventRequest{}
	api.Decode(r, &request)

	token := request.Token

	payload, err := json.Marshal(request.Data)
	if err != nil {
		return err
	}
	data := string(payload)

	switch request.Kind {
	case "js.error":
		e, err := models.Errors.Create(token, data)
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

		err = ctx.Queue.Push("create.error", payload)
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

		break
	case "js.log":
		logger.Info("js.log", "payload", data)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
