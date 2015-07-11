package models

import (
	"fmt"
	"strings"

	"github.com/erraroo/erraroo/logger"
)

// EventsStore is the interface to error data
type ErrorsStore interface {
	FindOrCreate(*Project, *Event) (*Error, error)
	FindQuery(ErrorQuery) (ErrorResults, error)
	FindByID(int64) (*Error, error)
	Update(*Error) error
	Touch(*Error) error
	AddTags(*Error, []Tag) error
}

type ErrorQuery struct {
	ProjectID int64
	Status    string
	Tags      []Tag
	Libaries  []int64
	QueryOptions
}

type ErrorResults struct {
	Errors []*Error
	Tags   []TagValue
	Total  int64
	Query  ErrorQuery
}

type errorsStore struct{ *Store }

func (s *errorsStore) FindOrCreate(p *Project, e *Event) (*Error, error) {
	er := newError(p, e)

	scope := s.Where(Error{
		Checksum:  er.Checksum,
		ProjectID: er.ProjectID,
	})

	if err := scope.Attrs(er).FirstOrCreate(&er).Error; err != nil {
		logger.Error("finding or creating error", "err", err)
		return nil, err
	}

	err := Libaries.Add(er, e.Libaries())
	if err != nil {
		logger.Error("adding Libaries", "err", err, "project", p.ID, "error.ID", er.ID)
		return nil, err
	}

	//err = Errors.Touch(er)
	//if err != nil {
	//logger.Error("touching e", "err", err, "e", e.ID)
	//return nil, err
	//}

	return er, nil
}

func (s *errorsStore) Touch(g *Error) error {
	query := "update errors set occurrences=(select count(*) from events where events.checksum = $1), last_seen_at=now_utc(), resolved='f', updated_at=now_utc() where id=$2 returning resolved, occurrences, last_seen_at, updated_at"
	return s.QueryRow(query, g.Checksum, g.ID).Scan(
		&g.Resolved,
		&g.Occurrences,
		&g.LastSeenAt,
		&g.UpdatedAt,
	)
}

func (s *errorsStore) Update(e *Error) error {
	return s.Save(e).Error
}

func (s *errorsStore) FindQuery(q ErrorQuery) (ErrorResults, error) {
	errors := ErrorResults{
		Query:  q,
		Errors: []*Error{},
	}

	scope := s.Table("errors").Debug()
	scope = scope.Where("errors.project_id=?", q.ProjectID)

	switch q.Status {
	case "unresolved":
		scope = scope.Where("errors.resolved=? and errors.muted=?", false, false)
	case "resolved":
		scope = scope.Where("errors.resolved=? and errors.muted=?", true, false)
	case "muted":
		scope = scope.Where("errors.muted=?", true)
	}

	joins := []string{}
	for i, lib := range q.Libaries {
		name := fmt.Sprintf("el_%d", i)

		join := fmt.Sprintf("inner join error_libraries %s on (%s.error_id=errors.id)", name, name)
		joins = append(joins, join)

		where := fmt.Sprintf("%s.library_id=?", name)
		scope = scope.Where(where, lib)
	}

	scope = scope.Joins(strings.Join(joins, " "))

	o := scope.Count(&errors.Total)
	if o.Error != nil {
		return errors, o.Error
	}

	scope = scope.Limit(q.PerPageOrDefault()).Offset(q.Offset())
	scope.Order("errors.last_seen_at desc")

	o = scope.Find(&errors.Errors)
	if o.Error != nil {
		return errors, o.Error
	}

	return errors, nil
}

func (s *errorsStore) FindByID(id int64) (*Error, error) {
	e := &Error{}
	o := s.First(&e, id)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	if o.Error != nil {
		return nil, o.Error
	}

	return e, nil
}

func (s *errorsStore) AddTags(e *Error, tags []Tag) error {
	for _, tag := range tags {
		update := "update project_keys set values_seen=values_seen+1 where project_id=$1 and key=$2"
		insert := "insert into project_keys (project_id, key, label) select $1,$2,$3"
		upsert := "with upsert as (%s returning *) %s where not exists (select * from upsert);"
		query := fmt.Sprintf(upsert, update, insert)

		_, err := s.Exec(query, e.ProjectID, tag.Key, tag.Label)
		if err != nil {
			logger.Error("inserting project_keys", "err", err)
			return err
		}

		update = "update project_key_values set times_seen=times_seen+1, last_seen=now_utc() where project_id=$1 and key=$2 and value=$3"
		insert = "insert into project_key_values (project_id, key, value) select $1,$2,$3"
		upsert = "with upsert as (%s returning *) %s where not exists (select * from upsert);"
		query = fmt.Sprintf(upsert, update, insert)

		_, err = s.Exec(query, e.ProjectID, tag.Key, tag.Value)
		if err != nil {
			logger.Error("inserting project_key_values", "err", err)
			return err
		}

		update = "update error_keys set values_seen=values_seen+1 where project_id=$1 and key=$2 and error_id=$3"
		insert = "insert into error_keys (project_id, key, error_id) select $1,$2,$3"
		upsert = "with upsert as (%s returning *) %s where not exists (select * from upsert);"
		query = fmt.Sprintf(upsert, update, insert)

		_, err = s.Exec(query, e.ProjectID, tag.Key, e.ID)
		if err != nil {
			logger.Error("inserting error_keys", "err", err)
			return err
		}

		update = "update error_key_values set times_seen=times_seen+1,last_seen=now_utc() where project_id=$1 and key=$2 and error_id=$3 and value=$4"
		insert = "insert into error_key_values (project_id, key, error_id, value) select $1,$2,$3,$4"
		upsert = "with upsert as (%s returning *) %s where not exists (select * from upsert);"
		query = fmt.Sprintf(upsert, update, insert)

		_, err = s.Exec(query, e.ProjectID, tag.Key, e.ID, tag.Value)
		if err != nil {
			logger.Error("inserting error_key_values", "err", err)
			return err
		}
	}

	return nil
}
