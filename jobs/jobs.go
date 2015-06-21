package jobs

import (
	"encoding/json"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/usecases"
	"github.com/nerdyworm/rsq"
)

var queue rsq.Queue

func Use(q rsq.Queue) {
	queue = q
}

func Work(handler rsq.JobHandler) {
	queue.Work(handler)
}

func Push(name string, payload []byte) error {
	return queue.Push(name, payload)
}

func Shutdown() error {
	return queue.Shutdown()
}

func EventProcess(job *rsq.Job, ctx *cx.Context) error {
	var id int64
	err := json.Unmarshal(job.Payload, &id)
	if err != nil {
		return err
	}

	return usecases.ErrorCreated(id)
}
