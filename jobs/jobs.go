package jobs

import "github.com/nerdyworm/rsq"

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
