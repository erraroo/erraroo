package usecases

import (
	"testing"

	"github.com/erraroo/erraroo/models"
	"github.com/stretchr/testify/assert"
)

func TestErrorCreatedCreatesGroup(t *testing.T) {
	_, _, project := aup(t)
	e := makeError(t, project, "{}")

	err := ErrorCreated(project, e)
	assert.Nil(t, err)

	groups, err := models.Groups.FindQuery(models.GroupQuery{ProjectID: project.ID})
	assert.Nil(t, err)
	assert.NotEmpty(t, groups)
}

func TestErrorCreated_DeliversNotifcations(t *testing.T) {
	emailSender.Clear()
	_, user, project := aup(t)
	e := makeError(t, project, "{}")

	err := ErrorCreated(project, e)
	assert.Nil(t, err)
	assert.Equal(t, len(emailSender.sends), 1)

	send := emailSender.sends[0]
	assert.Equal(t, send["to"], user.Email)

	err = ErrorCreated(project, e)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(emailSender.sends))
}

func TestErrorCreated_DeliversNotifcationsWhenResolved(t *testing.T) {
	emailSender.Clear()

	_, _, project := aup(t)
	e := makeError(t, project, "{}")

	err := ErrorCreated(project, e)
	assert.Equal(t, 1, len(emailSender.sends))

	group, err := models.Groups.FindOrCreate(project, e)
	group.Resolved = true
	err = models.Groups.Update(group)
	assert.Nil(t, err)

	emailSender.Clear()
	err = ErrorCreated(project, e)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(emailSender.sends))

	group, err = models.Groups.FindByID(group.ID)
	assert.Nil(t, err)
	assert.Equal(t, false, group.Resolved)

	group.Resolved = true
	group.Muted = true
	err = models.Groups.Update(group)
	assert.Nil(t, err)

	emailSender.Clear()
	err = ErrorCreated(project, e)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(emailSender.sends))
}
