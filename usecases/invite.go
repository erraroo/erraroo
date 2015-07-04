package usecases

import (
	"encoding/json"

	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
)

func InviteByEmail(from *models.User, to string) (*models.Invitation, error) {
	invitation, err := models.Invitations.Create(to, from)
	if err != nil {
		logger.Error("creating invitation", "err", err)
		return nil, err
	}

	payload, err := json.Marshal(invitation.Token)
	if err != nil {
		return invitation, err
	}

	err = jobs.Push("invitation.deliver", payload)
	if err != nil {
		return nil, err
	}

	return invitation, err
}

func InvitationDeliver(token string) error {
	invitation, err := models.Invitations.FindByToken(token)
	if err != nil {
		logger.Error("finding invitation", "token", token, "err", err)
		return err
	}

	err = mailers.DeliverInvitation(invitation)
	if err != nil {
		logger.Error("delivering invitation", "err", err)
		return err
	}

	return nil
}
