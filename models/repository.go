package models

type Repository struct {
	ID              int64
	ProjectID       int64
	Provider        string
	GithubOrg       string
	GithubScope     string
	GithubRepo      string
	GithubToken     string
	GithubTokenType string
}

func InsertRepository(r *Repository) error {
	return store.Save(r).Error
}

func FindRepositoryByID(id int64) (*Repository, error) {
	r := &Repository{}

	o := store.Where("id = ?", id).First(r)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	if o.Error != nil {
		return nil, o.Error
	}

	return r, nil
}

func FindRepositoryByProjectID(projectID int64) (*Repository, error) {
	r := &Repository{}

	o := store.Where("project_id = ?", projectID).First(r)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	if o.Error != nil {
		return nil, o.Error
	}

	return r, nil
}

func DeleteRepository(repository *Repository) error {
	return store.Delete(repository).Error
}

func FindRepositoryByGithubOrgAndGithubRepo(org string, repo string) (*Repository, error) {
	// webhooks come in with these two pieces of information we need to resolve
	return nil, nil
}
