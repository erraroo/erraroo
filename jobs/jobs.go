package jobs

import (
	"encoding/json"

	"github.com/erraroo/erraroo/logger"
	"github.com/nerdyworm/rsq"
)

var queue rsq.Queue

func Use(q rsq.Queue) {
	queue = q
}

func Work(handler rsq.JobHandler) {
	queue.Work(handler)
}

func Push(name string, payload interface{}) error {
	p, err := json.Marshal(payload)
	if err != nil {
		logger.Error("could not marshal payload", "err", err)
		return err
	}

	return queue.Push(name, p)
}

func Shutdown() error {
	return queue.Shutdown()
}
