package serializers

import "github.com/erraroo/erraroo/models"

type Repository struct {
	ID         int64
	ProjectID  int64
	GithubOrg  string
	GithubRepo string
}

type ShowRepository struct {
	Repository Repository
}

func NewShowRepository(r *models.Repository) ShowRepository {
	return ShowRepository{
		Repository: Repository{
			ID:         r.ID,
			ProjectID:  r.ProjectID,
			GithubOrg:  r.GithubOrg,
			GithubRepo: r.GithubRepo,
		},
	}
}
