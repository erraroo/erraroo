package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/erraroo/erraroo/logger"
)

type Revision struct {
	ID           int64
	ProjectID    int64
	SHA          string `json:"Sha"`
	Dependencies []Dependency
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (o *Revision) Empty() bool {
	return len(o.Dependencies) == 0
}

type Dependency struct {
	Latest   string `json:"latest"`
	Location string `json:"location"`
	Name     string `json:"name"`
	Target   string `json:"target"`
}

func (d Dependency) String() string {
	return fmt.Sprintf("[%s] %s %s -> %s", d.Location, d.Name, d.Target, d.Latest)
}

func FindRevisionsByProjectID(projectID int64) ([]*Revision, error) {
	revisions := []*Revision{}

	query := "select id, sha, dependencies from revisions where project_id = $1"
	rows, err := store.Query(query, projectID)
	if err != nil {
		logger.Error("query", "project", projectID, "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payload []byte
		revision := &Revision{
			ProjectID: projectID,
		}

		err := rows.Scan(
			&revision.ID,
			&revision.SHA,
			&payload,
		)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(payload, &revision.Dependencies)
		if err != nil {
			return nil, err
		}

		revisions = append(revisions, revision)
	}

	return revisions, err
}

func SaveRevision(revision *Revision) error {
	query := `WITH
	updated AS (
		update revisions set dependencies=$3, updated_at=now_utc() where project_id=$1 and sha=$2 RETURNING *
	),
	inserted AS (
		INSERT INTO revisions (project_id, sha, dependencies) SELECT $1,$2,$3 WHERE NOT EXISTS (SELECT 1 FROM updated WHERE project_id = $1 AND sha = $2) RETURNING *
	)
	SELECT id FROM inserted UNION ALL SELECT id from updated;`

	dependencies, err := json.Marshal(revision.Dependencies)
	if err != nil {
		logger.Error("marshaling dependencies", "project", revision.ProjectID, "err", err)
		return err
	}

	err = store.QueryRow(query, revision.ProjectID, revision.SHA, dependencies).Scan(
		&revision.ID,
	)

	if err != nil {
		logger.Error("inserting revision revision", "project", revision.ProjectID, "err", err)
		return err
	}

	return nil
}
