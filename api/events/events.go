package events

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
)

type CreateEventRequest struct {
	Kind  string                 `json:"kind"`
	Token string                 `json:"token"`
	Data  map[string]interface{} `json:"data"`
}

func Create(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := CreateEventRequest{}
	cx.Decode(r, &request)

	token := request.Token

	payload, err := json.Marshal(request.Data)
	if err != nil {
		return err
	}

	data := string(payload)
	log.Println(request)

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
	case "js.performance.timing":
		_, err := models.Timings.Create(token, data)
		if err == models.ErrNotFound {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		if err != nil {
			return err
		}

		break
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
