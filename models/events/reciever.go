package events

import (
	"encoding/json"

	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
)

type CreateEventRequest struct {
	Client struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"client"`

	Kind    string                 `json:"kind"`
	Data    map[string]interface{} `json:"data"`
	Session string                 `json:"session"`
}

func Ingest(token string, request CreateEventRequest) error {
	project, err := models.Projects.FindByToken(token)
	if err != nil {
		return err
	}

	switch request.Kind {
	case "js.error":
		err := jobs.Push("create.js.error", map[string]interface{}{
			"data":      request.Data,
			"projectID": project.ID,
		})

		if err != nil {
			return err
		}

	case "js.timing":
		raw, err := json.Marshal(request.Data)
		if err != nil {
			logger.Error("could not unmarshal timing data", "request", request, "err", err)
			return err
		}

		_, err = models.Timings.Create(project, string(raw))
		if err != nil {
			return err
		}
	}

	return nil
}
