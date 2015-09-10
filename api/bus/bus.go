package bus

import (
	"time"

	"github.com/erraroo/erraroo/logger"
	"github.com/nerdyworm/puller"
)

var Puller *puller.Puller

type Notifcation struct {
	Name    string
	Payload interface{}
}

func Push(channel string, payload interface{}) error {
	if Puller == nil {
		logger.Warn("pushing into nil puller", "channel", channel, "payload", payload)
		return nil
	}

	return Puller.Push(channel, payload)
}

func Pull(channels puller.Channels, t time.Duration) (puller.Backlog, error) {
	return Puller.Pull(channels, t)
}
