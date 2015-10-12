package usecases

import (
	"os"
	"testing"

	"github.com/erraroo/erraroo/models"
	"github.com/stretchr/testify/assert"
)

func TestOutdatedEmberProject_ThatDoesNotHaveARepositoriesitory(t *testing.T) {
	account := account(t)
	project := project(t, account)

	err := OutdatedEmberProject(project, nil)
	if err != ErrNoRepo {
		t.Fatalf("expected ErrNoRepo got %v", err)
	}
}

func TestOutdatedEmberProject_ThatDoesHasRepositoriesitory(t *testing.T) {
	account := account(t)
	project := project(t, account)

	models.InsertRepository(&models.Repository{
		ProjectID:   project.ID,
		Provider:    "github",
		GithubToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
		GithubOrg:   "erraroo",
		GithubRepo:  "erraroo-app",
	})

	checker := &githubNodeDepencyChecker{}
	err := OutdatedEmberProject(project, checker)
	assert.Nil(t, err)

	revisions, err := models.FindOutdatedRevisionsByProjectID(project.ID)
	assert.Nil(t, err)
	assert.Len(t, revisions, 1)
}
