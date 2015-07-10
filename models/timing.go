package models

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/erraroo/erraroo/logger"
)

type Timing struct {
	ID        int64
	CreatedAt time.Time `db:"created_at"`
	Payload   string
	ProjectID int64 `db:"project_id"`
}

type timingData map[string]float64

func (t *Timing) Average(payload string) error {
	newData := timingData{}
	err := json.Unmarshal([]byte(payload), &newData)
	if err != nil {
		log.Println(err)
		return err
	}

	oldData := timingData{}
	err = json.Unmarshal([]byte(t.Payload), &oldData)
	if err != nil {
		logger.Error("Unmarshal", "err", err)
		return err
	}

	for key := range newData {
		if newData[key] > 0 && oldData[key] > 0 {
			newData[key] = (oldData[key] + newData[key]) / 2
		} else if oldData[key] == 0 {
			newData[key] = newData[key]
		} else {
			newData[key] = 0
		}

		newData[key] = toFixed(newData[key], 0)
	}

	newPayload, err := json.Marshal(newData)
	if err != nil {
		logger.Error("Marshal", "err", err)
		return err
	}

	t.Payload = string(newPayload)
	return Timings.Update(t)
}

func (t *Timing) PreProcess() {
	data := timingData{}

	err := json.Unmarshal([]byte(t.Payload), &data)
	if err != nil {
		logger.Error("could not unmarshal timing payload", "err", err)
		return
	}

	for key := range data {
		data[key] = toFixed(data[key], 0)
	}

	payload, err := json.Marshal(data)
	if err != nil {
		logger.Error("could not marshal timing payload", "err", err, "data", data)
		return
	}

	t.Payload = string(payload)
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
