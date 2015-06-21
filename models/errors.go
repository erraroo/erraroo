package models

import (
	"database/sql"
	"time"

	"github.com/erraroo/erraroo/logger"
)

// EventsStore is the interface to error data
type ErrorsStore interface {
	FindOrCreate(*Project, *Event) (*Error, error)
	FindQuery(ErrorQuery) (ErrorResults, error)
	FindByID(int64) (*Error, error)
	Update(*Error) error
	Touch(*Error) error
}

type ErrorQuery struct {
	ProjectID int64
	Resolved  bool
	QueryOptions
}

type ErrorResults struct {
	Errors []*Error
	Total  int64
	Query  ErrorQuery
}
type errorsStore struct{ *Store }

func (s *errorsStore) FindOrCreate(p *Project, e *Event) (*Error, error) {
	group, err := s.findByProjectIDAndChecksum(p.ID, e.Checksum)
	if err == ErrNotFound {
		group = newError(p, e)
		return group, s.insert(group)
	} else if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *errorsStore) findByProjectIDAndChecksum(id int64, checksum string) (*Error, error) {
	group := &Error{}
	query := "select * from errors where project_id=$1 and checksum=$2 limit 1"
	err := s.Get(group, query, id, checksum)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *errorsStore) insert(group *Error) error {
	query := "insert into errors (name, message, checksum, project_id, occurrences, last_seen_at) values($1,$2,$3,$4,$5,$6) returning id"
	err := s.QueryRow(query,
		group.Name,
		group.Message,
		group.Checksum,
		group.ProjectID,
		group.Occurrences,
		group.LastSeenAt,
	).Scan(&group.ID)

	if err != nil {
		logger.Error("error inserting group", "err", err)
		return err
	}

	group.WasInserted = true
	return nil
}

func (s *errorsStore) Touch(g *Error) error {
	query := "update errors set occurrences=(select count(*) from events where events.checksum = $1), last_seen_at=now_utc(), resolved='f', updated_at=now_utc() where id=$2 returning occurrences, updated_at"
	return s.QueryRow(query, g.Checksum, g.ID).
		Scan(&g.Occurrences, &g.UpdatedAt)
}

func (s *errorsStore) Update(group *Error) error {
	group.UpdatedAt = time.Now().UTC()

	query := "update errors set occurrences=$1, last_seen_at=$2, resolved=$3, updated_at=$4, muted=$5 where id=$6"
	_, err := s.Exec(query,
		group.Occurrences,
		group.LastSeenAt,
		group.Resolved,
		group.UpdatedAt,
		group.Muted,
		group.ID,
	)

	return err
}

func (s *errorsStore) FindQuery(q ErrorQuery) (ErrorResults, error) {
	errors := ErrorResults{}
	errors.Query = q
	errors.Errors = []*Error{}

	countQuery := builder.Select("count(*)").From("errors")
	findQuery := builder.Select("*").From("errors")

	countQuery = countQuery.Where("project_id=?", q.ProjectID)
	findQuery = findQuery.Where("project_id=?", q.ProjectID)

	findQuery = findQuery.Limit(uint64(q.PerPageOrDefault())).Offset(uint64(q.Offset()))

	findQuery = findQuery.OrderBy("last_seen_at desc, created_at desc")

	query, args, _ := findQuery.ToSql()
	err := s.Select(&errors.Errors, query, args...)
	if err != nil {
		return errors, err
	}

	query, args, _ = countQuery.ToSql()
	err = s.Get(&errors.Total, query, args...)
	if err != nil {
		return errors, err
	}

	return errors, err
}

func (s *errorsStore) FindByID(id int64) (*Error, error) {
	group := &Error{}
	query := "select * from errors where id=$1 limit 1"
	err := s.Get(group, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return group, nil
}
