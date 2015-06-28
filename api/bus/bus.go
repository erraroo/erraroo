package bus

import (
	"time"

	"gopkg.in/redis.v3"

	"github.com/nerdyworm/puller"
)

var p *puller.Puller

func init() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	p = puller.New(puller.Options{
		MaxBacklogSize: 10,
		Redis:          client,
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
