package usecases

import (
	"encoding/json"
	"testing"

	"github.com/erraroo/erraroo/models"
	"github.com/nerdyworm/rsq"

	"github.com/stretchr/testify/assert"
)

func TestInviteByEmail_DeliversEmail(t *testing.T) {
	email := uniqEmail()

	_, user, _ := aup(t)
	invite, err := InviteByEmail(user, email)
	assert.Nil(t, err)
	assert.Equal(t, email, invite.Address)

	handler := rsq.NewJobRouter()
	handler.Handle("invitation.deliver", func(job *rsq.Job) error {
		var token string
		json.Unmarshal(job.Payload, &token)

		assert.Equal(t, token, invite.Token, "passed correct token into job")
		return InvitationDeliver(token)
	})

	queue.Work(handler)
	assert.Equal(t, 1, len(emailSender.sends))
	assert.Equal(t, email, emailSender.sends[0]["to"])

	invites, err := models.Invitations.ListForUser(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, invites)
	invite = invites[0]
	assert.Equal(t, email, invite.Address)
	assert.Contains(t, emailSender.sends[0]["body"], invite.Token)
}
