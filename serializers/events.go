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

func NewShowEvent(p *models.Event) ShowEvent {
	return ShowEvent{
		Event: Event{p},
	}
}

type Events struct {
	Events []Event
	Meta   struct {
		Pagination Pagination
	}
}

func NewEvents(results models.EventResults) Events {
	e := Events{}
	e.Events = make([]Event, len(results.Events))

	for i, p := range results.Events {
		e.Events[i] = Event{p}
	}

	pages := math.Ceil(float64(results.Total) / float64(results.Query.PerPageOrDefault()))
	e.Meta.Pagination.Pages = int(pages)
	e.Meta.Pagination.Page = results.Query.PageOrDefault()
	e.Meta.Pagination.Limit = results.Query.PerPageOrDefault()
	e.Meta.Pagination.Total = int(results.Total)

	return e
}
