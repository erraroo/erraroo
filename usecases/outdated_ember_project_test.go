package usecases

import (
	"testing"

	"github.com/erraroo/erraroo/models"
	"github.com/stretchr/testify/assert"
)

func TestCheckEmberDependencies_ThatDoesNotHaveARepositoriesitory(t *testing.T) {
	account := account(t)
	project := project(t, account)

	err := CheckEmberDependencies(project.ID, nil)
	if err != ErrNoRepo {
		t.Fatalf("expected ErrNoRepo got %v", err)
	}
}

type mockChecker struct{}

func (mock *mockChecker) Outdated(r *models.Repository) (*models.Revision, error) {
	revision := &models.Revision{
		ProjectID: r.ProjectID,
		Dependencies: []models.Dependency{
			models.Dependency{},
		},
	}

	return revision, nil
}

func TestCheckEmberDependencies_ThatDoesHasRepositoriesitory(t *testing.T) {
	account := account(t)
	project := project(t, account)
	repository(t, project)

	checker := &mockChecker{}
	err := CheckEmberDependencies(project.ID, checker)
	assert.Nil(t, err)

	revisions, err := models.FindRevisionsByProjectID(project.ID)
	assert.Nil(t, err)
	assert.Len(t, revisions, 1)
	assert.Len(t, revisions[0].Dependencies, 1)

	// Ensure it only inserts one per sha
	err = CheckEmberDependencies(project.ID, checker)
	assert.Nil(t, err)

	revisions, err = models.FindRevisionsByProjectID(project.ID)
	assert.Nil(t, err)
	assert.Len(t, revisions, 1)
}
