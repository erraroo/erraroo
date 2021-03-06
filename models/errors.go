package models

import "github.com/erraroo/erraroo/logger"

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
	Status    string
	Libaries  []int64
	QueryOptions
}

type ErrorResults struct {
	Errors []*Error
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

	scope := s.Table("errors")
	scope = scope.Where("errors.project_id=?", q.ProjectID)

	switch q.Status {
	case "unresolved":
		scope = scope.Where("errors.resolved=? and errors.muted=?", false, false)
	case "resolved":
		scope = scope.Where("errors.resolved=? and errors.muted=?", true, false)
	case "muted":
		scope = scope.Where("errors.muted=?", true)
	}

	o := scope.Count(&errors.Total)
	if o.Error != nil {
		return errors, o.Error
	}

	scope = scope.Limit(q.PerPageOrDefault()).Offset(q.Offset())
	scope = scope.Order("errors.last_seen_at desc")

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
