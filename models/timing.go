package models

import (
	"encoding/json"
	"log"
	"time"
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
		log.Println(err)
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
	}

	newPayload, err := json.Marshal(newData)
	if err != nil {
		log.Println(err)
		return err
	}

	t.Payload = string(newPayload)
	return Timings.Update(t)
}
