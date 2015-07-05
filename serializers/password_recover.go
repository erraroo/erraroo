package serializers

import "github.com/erraroo/erraroo/models"

type PasswordRecover struct {
	*models.PasswordRecover
	ID string
}

type ShowPasswordRecover struct {
	PasswordRecover PasswordRecover
}

func NewShowPasswordRecover(pr *models.PasswordRecover) ShowPasswordRecover {
	return ShowPasswordRecover{
		PasswordRecover: PasswordRecover{pr, pr.Token},
	}
}
