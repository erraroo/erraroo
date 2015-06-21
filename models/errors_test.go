package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	account, err := Accounts.Create()
	assert.Nil(t, err)

	project, err := Projects.Create("Test Project", account.ID)
	assert.Nil(t, err)

	e, err := Events.Create(project.Token, "{}")
	assert.Nil(t, err)

	group, err := Errors.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, group.Checksum, e.Checksum)
	assert.Equal(t, group.ProjectID, e.ProjectID)
	assert.Equal(t, group.Message, e.Message())
	assert.Equal(t, group.Name, e.Name())

	group2, err := Errors.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group2)
	assert.Equal(t, group.ID, group2.ID)

	groups, err := Errors.FindQuery(ErrorQuery{
		ProjectID: project.ID,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, groups.Errors)
	assert.True(t, groups.Total == 1)
	assert.Equal(t, groups.Errors[0].ID, group.ID)

	groups2, err := Errors.FindQuery(ErrorQuery{})
	assert.Nil(t, err)
	assert.Empty(t, groups2.Errors)

	groups, err = Errors.FindQuery(ErrorQuery{
		ProjectID: project.ID,
		QueryOptions: QueryOptions{
			PerPage: 1,
			Page:    1,
		},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, groups.Errors)

	groups, err = Errors.FindQuery(ErrorQuery{
		ProjectID: project.ID,
		QueryOptions: QueryOptions{
			PerPage: 1,
			Page:    2,
		},
	})
	assert.Nil(t, err)
	assert.Empty(t, groups.Errors)
}

func TestErrors_TouchCountsNumberOfOccurrences(t *testing.T) {
	account, err := Accounts.Create()
	assert.Nil(t, err)

	project, err := Projects.Create("Test Project", account.ID)
	assert.Nil(t, err)

	event, err := Events.Create(project.Token, "{}")
	assert.Nil(t, err)

	event, err = Events.Create(project.Token, "{}")
	assert.Nil(t, err)

	e, err := Errors.FindOrCreate(project, event)
	assert.Nil(t, err)

	err = Errors.Touch(e)
	assert.Nil(t, err)
	assert.Equal(t, e.Occurrences, 2)
}
