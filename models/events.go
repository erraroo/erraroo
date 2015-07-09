package models

import (
	"time"

	"github.com/erraroo/erraroo/logger"
)

// EventsStore is the interface to error data
type EventsStore interface {
	Create(token, kind, data string) (*Event, error)
	ListForProject(*Project) ([]*Event, error)
	FindByID(int64) (*Event, error)
	FindQuery(EventQuery) (EventResults, error)
	Update(*Event) error
}

type eventsStore struct{ *Store }

func (s *eventsStore) Create(token, kind, data string) (*Event, error) {
	var err error
	project, err := Projects.FindByToken(token)
	if err != nil {
		return nil, err
	}

	e := &Event{}
	e.Kind = kind
	e.Payload = data
	e.ProjectID = project.ID
	e.CreatedAt = time.Now().UTC()
	e.UpdatedAt = e.CreatedAt

	err = e.PreProcess()
	if err != nil {
		logger.Error("error pre processing event", "kind", kind, "payload", data, "token", token, "err", err)
		return nil, err
	}

	query := "insert into events (payload, project_id, checksum, kind, created_at, updated_at) values ($1,$2,$3,$4,$5,$6) returning id"
	row := s.QueryRow(query,
		e.Payload,
		e.ProjectID,
		e.Checksum,
		e.Kind,
		e.CreatedAt,
		e.UpdatedAt,
	)

	err = row.Scan(&e.ID)

	return e, err
}

func (s *eventsStore) ListForProject(p *Project) ([]*Event, error) {
	events := []*Event{}
	query := s.dbGorm.Where("project_id=?", p.ID).Order("created_at desc").Limit(100)
	return events, query.Find(&events).Error
}

func (s *eventsStore) FindByID(id int64) (*Event, error) {
	e := &Event{}
	o := s.dbGorm.First(&e, id)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *eventsStore) Update(e *Event) error {
	return s.dbGorm.Save(e).Error
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

	scope := s.dbGorm.Table("events")
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
