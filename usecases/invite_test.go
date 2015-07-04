package usecases

import (
	"testing"

	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"

	"github.com/stretchr/testify/assert"
)

func InviteByEmail(from *models.User, to string) error {
	invitation, err := models.Invitations.Create(to, from)
	if err != nil {
		logger.Error("creating invitation", "err", err)
		return err
	}

	err = mailers.DeliverInvitation(invitation)
	if err != nil {
		logger.Error("delivering invitation", "err", err)
		return err
	}

	return err
}

func TestInviteByEmail_DeliversEmail(t *testing.T) {
	email := uniqEmail()

	emailSender.Clear()
	_, user, _ := aup(t)
	err := InviteByEmail(user, email)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(emailSender.sends))
	assert.Equal(t, email, emailSender.sends[0]["to"])

	invites, err := models.Invitations.ListForUser(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, invites)

	invite := invites[0]
	assert.Equal(t, email, invite.Address)
	assert.Contains(t, emailSender.sends[0]["body"], invite.Token)
}
