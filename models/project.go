package models

import (
	"database/sql"
	"time"

	"github.com/tuvistavie/securerandom"
)

// Project is the thing we associate with errors
type Project struct {
	ID              int64
	Name            string
	Token           string
	AccountID       int64     `db:"account_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	UnresolvedCount int       `db:"unresolved_count"`
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

	project := &Project{}
	project.Name = name
	project.Token = token
	project.AccountID = accountID
	project.CreatedAt = time.Now().UTC()
	project.UpdatedAt = project.CreatedAt

	query := "insert into projects (name, token, account_id, created_at, updated_at) values ($1,$2,$3,$4,$5) returning id"
	row := s.QueryRow(query, project.Name, project.Token, project.AccountID, project.CreatedAt, project.UpdatedAt)
	return project, row.Scan(&project.ID)
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

	err := s.Get(project, "select * from projects where token = $1 limit 1", token)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return project, err
}

const unresolvedCount = "(select count(*) from errors where errors.project_id = projects.id and errors.resolved = 'f') as unresolved_count"

func (s *projectsStore) FindByID(id int64) (*Project, error) {
	project := &Project{}

	query := "select *, " + unresolvedCount + " from projects where projects.id = $1"
	err := s.Get(project, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return project, err
}

func (s *projectsStore) ByAccountID(id int64) ([]*Project, error) {
	projects := []*Project{}
	query := "select *, " + unresolvedCount + " from projects where projects.account_id = $1"
	return projects, s.Select(&projects, query, id)
}

func (s *projectsStore) Update(project *Project) error {
	project.UpdatedAt = time.Now().UTC()
	query := "update projects set name=$1, token=$2, updated_at=$3 where id=$4"
	_, err := s.Exec(query,
		project.Name,
		project.Token,
		project.UpdatedAt,
		project.ID,
	)

	return err
}
