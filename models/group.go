package models

import "time"

type Group struct {
	ID          int64
	Message     string
	Checksum    string
	Occurrences int
	Resolved    bool
	LastSeenAt  time.Time `db:"last_seen_at"`
	ProjectID   int64     `db:"project_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	WasInserted bool `db:"-"`
}

type GroupQueryResults struct {
	Groups []*Group
	Total  int64
}

func newGroup(p *Project, e *Error) *Group {
	return &Group{ProjectID: p.ID, Message: e.Message(), Checksum: e.Checksum}
}
