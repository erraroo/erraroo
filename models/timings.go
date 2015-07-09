package models

type TimingsStore interface {
	Create(project *Project, data string) (*Timing, error)
	Update(*Timing) error
	Last7Days(*Project) ([]*Timing, error)
}

type timingsStore struct{ *Store }

func (s *timingsStore) Create(project *Project, data string) (*Timing, error) {
	timing := &Timing{}
	o := s.Where("project_id=? and created_at=date_trunc('hour', now_utc())", project.ID).First(&timing)
	if o.RecordNotFound() {
		timing.ProjectID = project.ID
		timing.Payload = data
		timing.PreProcess()
		query := `insert into timings (
			project_id,
			payload,
			created_at
		) values ($1,$2, date_trunc('hour', now_utc())) returning id, created_at`

		row := s.QueryRow(query, timing.ProjectID, timing.Payload)
		return timing, row.Scan(&timing.ID, &timing.CreatedAt)
	} else if o.Error != nil {
		return nil, o.Error
	}

	return timing, timing.Average(data)
}

func (s *timingsStore) Update(t *Timing) error {
	return s.Save(t).Error
}

func (s *timingsStore) Last7Days(project *Project) ([]*Timing, error) {
	timings := []*Timing{}
	return timings, s.Where("project_id = ? and created_at > now_utc()::date - interval '7d'", project.ID).Find(&timings).Error
}
