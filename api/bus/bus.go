package bus

import (
	"time"

	"github.com/nerdyworm/puller"
)

var p puller.Puller

func init() {
	p = puller.NewMemoryPuller(puller.MemoryPullerOptions{
		MaxMessagesPerChannel: 100,
	})
}

type Notifcation struct {
	Name    string
	Payload interface{}
}

func Push(channel string, payload interface{}) error {
	return p.Push(channel, payload)
}

func Pull(channels puller.Channels, t time.Duration) (puller.Backlog, error) {
	return p.Pull(channels, t)
}
