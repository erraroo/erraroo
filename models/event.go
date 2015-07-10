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

func NewEvent(p *Project, kind string, payload string) *Event {
	return &Event{
		Kind:      kind,
		Payload:   payload,
		ProjectID: p.ID,
	}
}

func (e *Event) Message() string {
	return e.Handler().Message()
}

func (e *Event) Name() string {
	return e.Handler().Name()
}

func (e *Event) PreProcess() error {
	return e.Handler().PreProcess()
}

func (e *Event) PostProcess() error {
	return e.Handler().PostProcess()
}

func (e *Event) Libaries() []Library {
	return e.Handler().Libaries()
}

func (e *Event) Handler() EventHandler {
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

func HandlerFor(kind string, e *Event) EventHandler {
	switch kind {
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
	Message() string
	Name() string
	PreProcess() error
	PostProcess() error
	Libaries() []Library
}
