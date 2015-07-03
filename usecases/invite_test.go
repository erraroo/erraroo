package usecases

import (
	"testing"

	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"

	"github.com/stretchr/testify/assert"
)

func InviteByEmail(from *models.User, to string) error {
	err := mailers.DeliverInvitation(from, to)
	return err
}

func TestInviteByEmail_DeliversEmail(t *testing.T) {
	emailSender.Clear()
	_, user, _ := aup(t)
	err := InviteByEmail(user, "bob@example.com")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(emailSender.sends))

	// assert that there is a link
	// go to the link
	// test api can create user
}
