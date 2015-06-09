package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroups(t *testing.T) {
	account, err := Accounts.Create()
	assert.Nil(t, err)

	project, err := Projects.Create("Test Project", account.ID)
	assert.Nil(t, err)

	e, err := Errors.Create(project.Token, "{}")
	assert.Nil(t, err)

	group, err := Groups.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, group.Checksum, e.Checksum)
	assert.Equal(t, group.ProjectID, e.ProjectID)
	assert.Equal(t, group.Message, e.Message())

	group2, err := Groups.FindOrCreate(project, e)
	assert.Nil(t, err)
	assert.NotNil(t, group2)
	assert.Equal(t, group.ID, group2.ID)

	groups, err := Groups.FindQuery(GroupQuery{
		ProjectID: project.ID,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, groups.Groups)
	assert.True(t, groups.Total == 1)
	assert.Equal(t, groups.Groups[0].ID, group.ID)

	groups2, err := Groups.FindQuery(GroupQuery{})
	assert.Nil(t, err)
	assert.Empty(t, groups2.Groups)

	groups, err = Groups.FindQuery(GroupQuery{
		ProjectID: project.ID,
		QueryOptions: QueryOptions{
			PerPage: 1,
			Page:    1,
		},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, groups.Groups)

	groups, err = Groups.FindQuery(GroupQuery{
		ProjectID: project.ID,
		QueryOptions: QueryOptions{
			PerPage: 1,
			Page:    2,
		},
	})
	assert.Nil(t, err)
	assert.Empty(t, groups.Groups)
}
