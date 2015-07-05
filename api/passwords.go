package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/gorilla/mux"
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

type CreatePasswordRecoverRequest struct {
	Email string
}

func PasswordRecoversCreate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := CreatePasswordRecoverRequest{}
	Decode(r, &request)

	user, err := models.Users.FindByEmail(request.Email)
	if err == models.ErrNotFound {
		errs := models.NewValidationErrors()
		errs.Add("Email", "not found")
		return errs
	}

	pr, err := models.PasswordRecovers.Create(user)
	if err != nil {
		logger.Error("creating password recover", "err", err)
		return err
	}

	err = mailers.DeliverPasswordRecover(pr)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusCreated, nil)
}

func PasswordRecoversShow(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	token := mux.Vars(r)["token"]
	logger.Debug(token)

	pr, err := models.PasswordRecovers.FindByToken(token)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowPasswordRecover(pr))
}

type PasswordRecoversUpdateRequest struct {
	Password string
	Token    string
	pr       *models.PasswordRecover
}

func (p *PasswordRecoversUpdateRequest) Validate() (models.ValidationErrors, error) {
	var err error
	errs := models.NewValidationErrors()

	if p.Password == "" {
		errs.Add("Password", "can't be blank")
	}

	if p.Token == "" {
		errs.Add("Token", "can't be blank")
	} else {
		p.pr, err = models.PasswordRecovers.FindByToken(p.Token)
		if err != nil {
			return errs, err
		}

		if p.pr.Used {
			errs.Add("Token", "has already been used")
		}
	}

	return errs, nil

}

func PasswordRecoversUpdate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := &PasswordRecoversUpdateRequest{}
	Decode(r, request)

	errors, err := request.Validate()
	if err != nil {
		return err
	}

	if errors.Any() {
		return errors
	}

	err = models.PasswordRecovers.Use(request.pr, request.Password)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, struct {
		Message string
	}{"Reset your password! Try logging in now"})
}
