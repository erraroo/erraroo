package models

import "time"

type ProjectKey struct {
	ProjectID  int64 `db:"project_id"`
	Key        string
	ValuesSeen int64
	Label      string
}

type ProjectKeyValue struct {
	ProjectID int64 `db:"project_id"`
	Key       string
	Value     string
	TimesSeen int64
	FirstSeen time.Time `db:"first_seen"`
	LastSeen  time.Time `db:"last_seen"`
}

type ErrorKey struct {
	ProjectID  int64 `db:"project_id"`
	Key        string
	ErrorID    int64 `db:"error_id"`
	ValuesSeen int64
}

type ErrorKeyValue struct {
	ProjectID int64 `db:"project_id"`
	Key       string
	Value     string
	ErrorID   int64 `db:"error_id"`
	TimesSeen int64
	FirstSeen time.Time `db:"first_seen"`
	LastSeen  time.Time `db:"last_seen"`
}
