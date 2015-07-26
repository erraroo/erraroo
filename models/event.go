package models

import (
	"fmt"
	"time"
)

// Event is the entity that stores error data
type Event struct {
	ID        int64
	Payload   string
	Checksum  string
	Kind      string
	ProjectID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEvent(p *Project, kind string, payload string) *Event {
	return &Event{
		Kind:      kind,
		Payload:   payload,
		ProjectID: p.ID,
	}
}

func (e *Event) BeforeCreate() error {
	return e.Handler().PreCreate()
}

func (e *Event) Message() string {
	return e.Handler().Message()
}

func (e *Event) Name() string {
	return e.Handler().Name()
}

func (e *Event) Libaries() []Library {
	return e.Handler().Libaries()
}

func (e *Event) Handler() EventHandler {
	switch e.Kind {
	case "js.error":
		return &jsErrorEvent{e}
	}

	return nil
}

func (e *Event) PayloadKey() string {
	t := e.CreatedAt.Format("2006/01/02")
	return fmt.Sprintf("projects/%d/events/%s/%d/payload.json", e.ProjectID, t, e.ID)
}

func (e *Event) SignedPayloadURL() string {
	return Events.PayloadURL(e)
}

type EventHandler interface {
	Checksum() string
	Message() string
	Name() string
	PreCreate() error
	Libaries() []Library
}
