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
