package models

import (
	"database/sql"
	"log"
)

type TimingsStore interface {
	Create(project *Project, data string) (*Timing, error)
	Update(*Timing) error
	Last7Days(*Project) ([]*Timing, error)
}

type timingsStore struct{ *Store }

func (s *timingsStore) Create(project *Project, data string) (*Timing, error) {
	timing := &Timing{}
	err := s.Get(timing, "select * from timings where project_id=$1 and created_at=date_trunc('minute', now());",
		project.ID,
	)

	if err == sql.ErrNoRows {
		timing.ProjectID = project.ID
		timing.Payload = data

		query := `insert into timings (
			project_id,
			payload,
			created_at
		) values ($1,$2, date_trunc('minute', now())) returning id, created_at`

		row := s.QueryRow(query, timing.ProjectID, timing.Payload)
		return timing, row.Scan(&timing.ID, &timing.CreatedAt)
	} else if err != nil {
		return nil, err
	}

	return timing, timing.Average(data)
}

func (s *timingsStore) Update(t *Timing) error {
	query := "update timings set payload=$1 where id=$2"
	_, err := s.Exec(query, t.Payload, t.ID)
	return err
}

func (s *timingsStore) Last7Days(project *Project) ([]*Timing, error) {
	timings := []*Timing{}

	query := "select * from timings where project_id = $1"
	err := s.Select(&timings, query, project.ID)
	if err != nil {
		log.Printf("[store] [error]: %v\n", err)
		return nil, err
	}

	return timings, nil
}
