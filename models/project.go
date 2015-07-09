package models

import (
	"fmt"
	"time"

	"github.com/tuvistavie/securerandom"
)

// Project is the thing we associate with errors
type Project struct {
	ID              int64
	Name            string
	Token           string
	AccountID       int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UnresolvedCount int `sql:"-"`
}

func (p *Project) Channel() string {
	return fmt.Sprintf("accounts.%d", p.AccountID)
}

// ProjectQuery is the options availabe to project searchs
type ProjectQuery struct {
	UserID int64
}

// ProjectsStore is the interface to project data
type ProjectsStore interface {
	Create(name string, accountID int64) (*Project, error)

	// GenerateToken returns a unqie token that can be used to identify a project
	GenerateToken() (string, error)

	// FindByToken returns the project for that token
	FindByToken(string) (*Project, error)
	FindByID(int64) (*Project, error)

	// ByAccountID returns the projects matching the account
	ByAccountID(int64) ([]*Project, error)

	// Update updates the project
	Update(*Project) error
}

type projectsStore struct{ *Store }

func (s *projectsStore) Create(name string, accountID int64) (*Project, error) {
	token, err := s.GenerateToken()
	if err != nil {
		return nil, err
	}

	project := &Project{
		Name:      name,
		Token:     token,
		AccountID: accountID,
	}

	return project, s.Save(project).Error
}

func (s *projectsStore) GenerateToken() (string, error) {
	for {
		token, err := securerandom.UrlSafeBase64(16, false)
		if err != nil {
			return "", err
		}

		exists := false
		row := s.QueryRow("select exists(select 1 from projects where token = $1)", token)
		err = row.Scan(&exists)
		if err != nil {
			return "", err
		}

		if !exists {
			return token, err
		}
	}
}

func (s *projectsStore) FindByToken(token string) (*Project, error) {
	project := &Project{}

	o := s.Where("token=?", token).First(&project)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	return project, o.Error
}

const unresolvedCount = "(select count(*) from errors where errors.project_id = projects.id and errors.resolved = 'f' and errors.muted = 'f') as unresolved_count"

func (s *projectsStore) FindByID(id int64) (*Project, error) {
	project := &Project{}
	sel := []string{"projects.*", unresolvedCount}
	out := s.Where("projects.id=?", id).Select(sel).Find(&project)
	if out.RecordNotFound() {
		return nil, ErrNotFound
	}

	return project, out.Error
}

func (s *projectsStore) ByAccountID(id int64) ([]*Project, error) {
	projects := []*Project{}
	sel := []string{"projects.*", unresolvedCount}
	out := s.Where("projects.account_id=?", id).Select(sel).Find(&projects)
	return projects, out.Error
}

func (s *projectsStore) Update(project *Project) error {
	return s.Save(project).Error
}
