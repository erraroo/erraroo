package serializers

import (
	"math"

	"github.com/erraroo/erraroo/models"
)

type Event struct {
	*models.Event
}

type ShowEvent struct {
	Event Event
}

func NewShowEvent(e *models.Event) ShowEvent {
	return ShowEvent{
		Event: NewEvent(e),
	}
}

type Events struct {
	Events []Event
	Meta   struct {
		Pagination Pagination
	}
}

func NewEvent(e *models.Event) Event {
	event := Event{
		Event: e,
	}

	return event
}

func NewEvents(results models.EventResults) Events {
	e := Events{}
	e.Events = make([]Event, len(results.Events))

	for i, event := range results.Events {
		e.Events[i] = NewEvent(event)
	}

	pages := math.Ceil(float64(results.Total) / float64(results.Query.PerPageOrDefault()))
	e.Meta.Pagination.Pages = int(pages)
	e.Meta.Pagination.Page = results.Query.PageOrDefault()
	e.Meta.Pagination.Limit = results.Query.PerPageOrDefault()
	e.Meta.Pagination.Total = int(results.Total)

	return e
}
