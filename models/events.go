package models

import "github.com/erraroo/erraroo/jobs"

// EventsStore is the interface to error data
type EventsStore interface {
	Create(token, kind, data string) (*Event, error)
	ListForProject(*Project) ([]*Event, error)
	FindByID(int64) (*Event, error)
	FindQuery(EventQuery) (EventResults, error)
	Update(*Event) error
	Insert(*Event) error
}

type eventsStore struct{ *Store }

func (s *eventsStore) Create(token, kind, data string) (*Event, error) {
	var err error
	project, err := Projects.FindByToken(token)
	if err != nil {
		return nil, err
	}

	switch kind {
	case "js.error":
		err = jobs.Push("create.js.error", map[string]string{
			"token": token,
			"data":  data,
		})

		if err != nil {
			return nil, err
		}

	case "js.timing":
		_, err := Timings.Create(project, data)
		if err != nil {
			return nil, err
		}

	}

	return nil, err
}

func (s *eventsStore) ListForProject(p *Project) ([]*Event, error) {
	events := []*Event{}
	query := s.Where("project_id=?", p.ID).Order("created_at desc").Limit(100)
	return events, query.Find(&events).Error
}

func (s *eventsStore) FindByID(id int64) (*Event, error) {
	e := &Event{}
	o := s.First(&e, id)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *eventsStore) Update(e *Event) error {
	err := s.Save(e).Error
	if err != nil {
		return err
	}

	return err
}

func (s *eventsStore) Insert(e *Event) error {
	err := e.PreProcess()
	if err != nil {
		return err
	}

	err = s.Save(e).Error
	if err != nil {
		return err
	}

	return err
}

type EventQuery struct {
	ProjectID int64
	Checksum  string
	Kind      string
	QueryOptions
}

type EventResults struct {
	Events []*Event
	Total  int64
	Query  EventQuery
}

func (s *eventsStore) FindQuery(q EventQuery) (EventResults, error) {
	events := EventResults{
		Query:  q,
		Events: []*Event{},
	}

	scope := s.Table("events")
	scope = scope.Where("events.project_id=?", q.ProjectID)

	if q.Checksum != "" {
		scope = scope.Where("events.checksum=?", q.Checksum)
	}

	if q.Kind != "" {
		scope = scope.Where("events.kind=?", q.Kind)
	}

	o := scope.Count(&events.Total)
	if o.Error != nil {
		return events, o.Error
	}

	scope = scope.Limit(q.PerPageOrDefault()).Offset(q.Offset())
	scope = scope.Order("events.created_at desc")

	o = scope.Find(&events.Events)
	if o.Error != nil {
		return events, o.Error
	}

	return events, nil
}
