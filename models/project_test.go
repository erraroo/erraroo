package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectsByAccountID(t *testing.T) {
	account1, err := Accounts.Create()
	assert.Nil(t, err)

	account2, err := Accounts.Create()
	assert.Nil(t, err)

	project1, err := Projects.Create("test project 1", account1.ID)
	assert.Nil(t, err)

	project2, err := Projects.Create("test project 2", account2.ID)
	assert.Nil(t, err)

	projects, err := Projects.ByAccountID(account1.ID)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(projects))
	assert.Equal(t, project1.ID, projects[0].ID)
	assert.Equal(t, project1.Name, projects[0].Name)
	assert.Equal(t, project1.AccountID, projects[0].AccountID)

	projects, err = Projects.ByAccountID(account2.ID)
	assert.Nil(t, err)
	assert.Equal(t, project2.ID, projects[0].ID)
	assert.Equal(t, project2.Name, projects[0].Name)
	assert.Equal(t, project2.AccountID, projects[0].AccountID)
}
