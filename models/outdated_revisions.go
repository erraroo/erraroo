package models

import (
	"encoding/json"

	"github.com/erraroo/erraroo/logger"
)

func FindOutdatedRevisionsByProjectID(projectID int64) ([]*OutdatedRevision, error) {
	revisions := []*OutdatedRevision{}

	query := "select id, sha, dependencies from outdated_revisions where project_id = $1"
	rows, err := store.Query(query, projectID)
	if err != nil {
		logger.Error("query", "project", projectID, "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payload []byte
		revision := &OutdatedRevision{
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

func InsertOutdatedRevision(outdated *OutdatedRevision) error {
	query := "insert into outdated_revisions (project_id, sha, dependencies) values($1,$2,$3) returning id"

	dependencies, err := json.Marshal(outdated.Dependencies)
	if err != nil {
		logger.Error("marshaling dependencies", "project", outdated.ProjectID, "err", err)
		return err
	}

	err = store.QueryRow(query, outdated.ProjectID, outdated.SHA, dependencies).Scan(
		&outdated.ID,
	)

	if err != nil {
		logger.Error("inserting outdated revision", "project", outdated.ProjectID, "err", err)
		return err
	}

	return nil
}
