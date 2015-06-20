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

	WasInserted bool `db:"-"`
}

type ErrorQueryResults struct {
	Errors []*Error
	Total  int64
}

func newError(p *Project, e *Event) *Error {
	return &Error{ProjectID: p.ID, Message: e.Message(), Checksum: e.Checksum}
}

func (g *Error) ShouldNotify() bool {
	if g.WasInserted {
		return true
	}

	if g.Resolved && !g.Muted {
		return true
	}

	return false
}
