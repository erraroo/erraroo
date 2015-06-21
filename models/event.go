package models

import "time"

// Event is the entity that stores error data
type Event struct {
	ID        int64
	Payload   string
	Checksum  string
	Kind      string
	ProjectID int64     `db:"project_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e *Event) Message() string {
	return e.handler().Message()
}

func (e *Event) Name() string {
	return e.handler().Name()
}

func (e *Event) IsAsync() bool {
	return e.handler().IsAsync()
}

func (e *Event) PreProcess() error {
	return e.handler().PreProcess()
}

func (e *Event) PostProcess() error {
	return e.handler().PostProcess()
}

func (e *Event) handler() EventHandler {
	switch e.Kind {
	case "js.error":
		return &jsErrorEvent{e}
	case "js.timing":
		return &jsTimingEvent{e}
	case "js.log":
		return &jsLogEvent{e}
	}

	return nil
}

type EventHandler interface {
	Checksum() string
	IsAsync() bool
	Message() string
	Name() string
	PreProcess() error
	PostProcess() error
}
