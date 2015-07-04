package serializers

import "github.com/erraroo/erraroo/models"

type Invitation struct {
	*models.Invitation
	ID string
}
type ShowInvitation struct {
	Invitation Invitation
}

func NewShowInvitation(i *models.Invitation) ShowInvitation {
	return ShowInvitation{
		Invitation: Invitation{i, i.Token},
	}
}

type Invitations struct {
	Invitations []Invitation
}

func NewInvitations(invitations []*models.Invitation) Invitations {
	payload := Invitations{}
	payload.Invitations = make([]Invitation, len(invitations))

	for i, ii := range invitations {
		payload.Invitations[i] = Invitation{ii, ii.Token}
	}

	return payload
}
