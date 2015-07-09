package models

import "time"

type Error struct {
	ID          int64
	Name        string
	Message     string
	Checksum    string
	Occurrences int
	Resolved    bool
	Muted       bool
	LastSeenAt  time.Time `db:"last_seen_at"`
	ProjectID   int64     `db:"project_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	Tags []TagValue `db:"-" json:"-" sql:"-"`
}

type ErrorQueryResults struct {
	Errors []*Error
	Total  int64
}

func newError(p *Project, e *Event) *Error {
	return &Error{
		Checksum:   e.Checksum,
		Name:       e.Name(),
		Message:    e.Message(),
		ProjectID:  p.ID,
		LastSeenAt: time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
	}
}

func (e *Error) ShouldNotify() bool {
	if e.Occurrences == 0 {
		return true
	}

	if e.Resolved && !e.Muted {
		return true
	}

	return false
}
