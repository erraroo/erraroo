package models

import "time"

type Library struct {
	ID        int64
	ProjectID int64
	Name      string
	Version   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ErrorLibrary struct {
	ErrorID   int64
	LibraryID int64
}

type LibariesStore interface {
	Add(*Error, []Library) error
	ListForProject(*Project) ([]*Library, error)
}

type libariesStore struct {
	*Store
}

func (s *libariesStore) Add(e *Error, libs []Library) error {
	for _, library := range libs {
		scope := s.Where(library)
		if err := scope.FirstOrCreate(&library).Error; err != nil {
			return err
		}

		el := new(ErrorLibrary)
		if err := s.Where(ErrorLibrary{e.ID, library.ID}).FirstOrCreate(&el).Error; err != nil {
			return err
		}

	}
	return nil
}

func (s *libariesStore) ListForProject(p *Project) ([]*Library, error) {
	libs := []*Library{}
	scope := s.Where("project_id=?", p.ID)
	return libs, scope.Find(&libs).Error
}
