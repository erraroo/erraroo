package bus

import (
	"time"

	"github.com/nerdyworm/puller"
)

var Puller *puller.Puller

type Notifcation struct {
	Name    string
	Payload interface{}
}

func Push(channel string, payload interface{}) error {
	return Puller.Push(channel, payload)
}

func Pull(channels puller.Channels, t time.Duration) (puller.Backlog, error) {
	return Puller.Pull(channels, t)
}
