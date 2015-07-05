package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
)

// ChangePasswordRequest the struct representing signup attributes
type ChangePasswordRequest struct {
	CurrentPassword    string
	NewPassword        string
	ConfirmNewPassword string
}

func (c ChangePasswordRequest) Validate() (models.ValidationErrors, error) {
	var err error

	errs := models.NewValidationErrors()
	if c.CurrentPassword == "" {
		errs.Add("CurrentPassword", "can't be blank")
	}

	if c.NewPassword == "" {
		errs.Add("NewPassword", "can't be blank")
	}

	if c.ConfirmNewPassword == "" {
		errs.Add("ConfirmNewPassword", "can't be blank")
	}

	if c.NewPassword != c.ConfirmNewPassword {
		errs.Add("ConfirmNewPassword", "does not match")
	}

	return errs, err
}

func PasswordsCreate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := ChangePasswordRequest{}
	Decode(r, &request)

	errs, err := request.Validate()
	if err != nil {
		return err
	}

	if errs.Any() {
		return errs
	}

	user, err := models.Users.FindByEmail(ctx.User.Email)
	if err != nil {
		return err
	}

	err = user.CheckPassword(request.CurrentPassword)
	if err != nil {
		errs.Add("CurrentPassword", "Your current password is not valid")
		return errs
	}

	user.SetPassword(request.NewPassword)

	err = models.Users.Update(user)
	if err != nil {
		logger.Error("updating user", "err", err)
		return err
	}

	return JSON(w, http.StatusOK, struct {
		Message string
	}{"Password changed"})
}
