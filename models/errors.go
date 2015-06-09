package models

import (
	"log"
	"time"
)

// ErrorsStore is the interface to error data
type ErrorsStore interface {
	Create(token, data string) (*Error, error)
	ListForProject(*Project) ([]*Error, error)
	FindByID(int64) (*Error, error)
	FindQuery(ErrorQuery) (ErrorResults, error)
	Update(*Error) error
}

type errorsStore struct{ *Store }

func (s *errorsStore) Create(token, data string) (*Error, error) {
	project, err := Projects.FindByToken(token)
	if err != nil {
		return nil, err
	}

	e := &Error{}
	e.Payload = data
	e.ProjectID = project.ID
	e.CreatedAt = time.Now().UTC()
	e.UpdatedAt = e.CreatedAt
	e.generateChecksum()

	query := "insert into errors (payload, project_id, checksum, created_at, updated_at) values ($1,$2,$3,$4,$5) returning id"
	row := s.QueryRow(query,
		e.Payload,
		e.ProjectID,
		e.Checksum,
		e.CreatedAt,
		e.UpdatedAt,
	)

	return e, row.Scan(&e.ID)
}

func (s *errorsStore) ListForProject(p *Project) ([]*Error, error) {
	query := "select * from errors where project_id = $1 order by created_at desc limit 100"
	errors := []*Error{}
	err := s.Select(&errors, query, p.ID)
	return errors, err
}

func (s *errorsStore) FindByID(id int64) (*Error, error) {
	e := &Error{}
	query := "select * from errors where id = $1 limit 1"
	return e, s.Get(e, query, id)
}

func (s *errorsStore) Update(e *Error) error {
	query := "update errors set payload=$1 where id = $2"

	_, err := s.Exec(query, e.Payload, e.ID)
	if err != nil {
		log.Printf("error updating error %v\n", err)
		return err
	}

	return nil
}

type ErrorQuery struct {
	ProjectID int64
	Checksum  string
	QueryOptions
}

type ErrorResults struct {
	Errors []*Error
	Total  int64
	Query  ErrorQuery
}

func (s *errorsStore) FindQuery(q ErrorQuery) (ErrorResults, error) {
	errs := ErrorResults{}
	errs.Query = q
	errs.Errors = []*Error{}

	countQuery := builder.Select("count(*)").From("errors")
	findQuery := builder.Select("*").From("errors")

	countQuery = countQuery.Where("project_id=?", q.ProjectID)
	findQuery = findQuery.Where("project_id=?", q.ProjectID)

	if q.Checksum != "" {
		countQuery = countQuery.Where("checksum=?", q.Checksum)
		findQuery = findQuery.Where("checksum=?", q.Checksum)
	}

	findQuery = findQuery.Limit(uint64(q.PerPageOrDefault())).Offset(uint64(q.Offset()))
	findQuery = findQuery.OrderBy("created_at desc")

	query, args, _ := findQuery.ToSql()
	err := s.Select(&errs.Errors, query, args...)
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
