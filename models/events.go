package models

import (
	"log"
	"time"
)

// EventsStore is the interface to error data
type EventsStore interface {
	Create(token, data string) (*Event, error)
	ListForProject(*Project) ([]*Event, error)
	FindByID(int64) (*Event, error)
	FindQuery(EventQuery) (EventResults, error)
	Update(*Event) error
}

type eventsStore struct{ *Store }

func (s *eventsStore) Create(token, data string) (*Event, error) {
	project, err := Projects.FindByToken(token)
	if err != nil {
		return nil, err
	}

	e := &Event{}
	e.Payload = data
	e.ProjectID = project.ID
	e.CreatedAt = time.Now().UTC()
	e.UpdatedAt = e.CreatedAt
	e.generateChecksum()

	query := "insert into events (payload, project_id, checksum, kind, created_at, updated_at) values ($1,$2,$3,$4,$5,$6) returning id"
	row := s.QueryRow(query,
		e.Payload,
		e.ProjectID,
		e.Checksum,
		e.Kind,
		e.CreatedAt,
		e.UpdatedAt,
	)

	return e, row.Scan(&e.ID)
}

func (s *eventsStore) ListForProject(p *Project) ([]*Event, error) {
	query := "select * from events where project_id = $1 order by created_at desc limit 100"
	events := []*Event{}
	err := s.Select(&events, query, p.ID)
	return events, err
}

func (s *eventsStore) FindByID(id int64) (*Event, error) {
	e := &Event{}
	query := "select * from events where id = $1 limit 1"
	return e, s.Get(e, query, id)
}

func (s *eventsStore) Update(e *Event) error {
	query := "update events set payload=$1 where id = $2"

	_, err := s.Exec(query, e.Payload, e.ID)
	if err != nil {
		log.Printf("error updating error %v\n", err)
		return err
	}

	return nil
}

type EventQuery struct {
	ProjectID int64
	Checksum  string
	QueryOptions
}

type EventResults struct {
	Events []*Event
	Total  int64
	Query  EventQuery
}

func (s *eventsStore) FindQuery(q EventQuery) (EventResults, error) {
	errs := EventResults{}
	errs.Query = q
	errs.Events = []*Event{}

	countQuery := builder.Select("count(*)").From("events")
	findQuery := builder.Select("*").From("events")

	countQuery = countQuery.Where("project_id=?", q.ProjectID)
	findQuery = findQuery.Where("project_id=?", q.ProjectID)

	if q.Checksum != "" {
		countQuery = countQuery.Where("checksum=?", q.Checksum)
		findQuery = findQuery.Where("checksum=?", q.Checksum)
	}

	findQuery = findQuery.Limit(uint64(q.PerPageOrDefault())).Offset(uint64(q.Offset()))
	findQuery = findQuery.OrderBy("created_at desc")

	query, args, _ := findQuery.ToSql()
	err := s.Select(&errs.Events, query, args...)
	if err != nil {
		return errs, err
	}

	query, args, _ = countQuery.ToSql()
	err = s.Get(&errs.Total, query, args...)
	if err != nil {
		return errs, err
	}

	return errs, err
}
