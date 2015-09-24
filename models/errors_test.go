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

	e := NewEvent(project, "js.error", "{}")
	err = Events.Insert(e)
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

func TestErrorsAreOrderedCorrectly(t *testing.T) {
	account, err := Accounts.Create()
	assert.Nil(t, err)

	project, err := Projects.Create("Test Project", account.ID)
	assert.Nil(t, err)

	event1 := NewEvent(project, "js.error", `{"trace":{"name":"event1"}}`)
	err = Events.Insert(event1)
	assert.Nil(t, err)

	error1, err := Errors.FindOrCreate(project, event1)
	assert.Nil(t, err)

	event2 := NewEvent(project, "js.error", `{"trace":{"name":"event2"}}`)
	err = Events.Insert(event2)
	assert.Nil(t, err)

	error2, err := Errors.FindOrCreate(project, event2)
	assert.Nil(t, err)

	query := ErrorQuery{ProjectID: project.ID}
	errors, err := Errors.FindQuery(query)
	assert.Nil(t, err)
	assert.NotEmpty(t, errors.Errors)
	assert.Equal(t, errors.Errors[0].ID, error1.ID)
	assert.Equal(t, errors.Errors[1].ID, error2.ID)

	Errors.Touch(error2)
	errors, err = Errors.FindQuery(query)
	assert.Nil(t, err)
	assert.NotEmpty(t, errors.Errors)
	assert.Equal(t, errors.Errors[0].ID, error2.ID)
	assert.Equal(t, errors.Errors[1].ID, error1.ID)
}

func TestErrors_TouchCountsNumberOfOccurrences(t *testing.T) {
	account, err := Accounts.Create()
	assert.Nil(t, err)

	project, err := Projects.Create("Test Project", account.ID)
	assert.Nil(t, err)

	event := NewEvent(project, "js.error", "{}")
	err = Events.Insert(event)
	assert.Nil(t, err)

	event = NewEvent(project, "js.error", "{}")
	err = Events.Insert(event)
	assert.Nil(t, err)

	e, err := Errors.FindOrCreate(project, event)
	assert.Nil(t, err)

	err = Errors.Touch(e)
	assert.Nil(t, err)
	assert.Equal(t, e.Occurrences, 2)
}
