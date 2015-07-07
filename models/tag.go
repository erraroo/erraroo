package models

import "time"

type Tag struct {
	Key   string
	Value string
	Label string
}

type TagValue struct {
	ID          int64
	ProjectID   int64 `db:"project_id"`
	ErrorID     int64 `db:"error_id"`
	Key         string
	Value       string
	Occurrences int64
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
