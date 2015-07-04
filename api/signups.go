package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

// SignupRequest an incoming sign up request
type SignupRequest struct {
	Signup Signup
}

// Validate the request to ensure that it is acceptable
func (s SignupRequest) Validate() (models.ValidationErrors, error) {
	errs := models.NewValidationErrors()
	if s.Signup.Email == "" {
		errs.Add("Email", "can't be blank")
	} else if models.Users.Exists(s.Signup.Email) {
		errs.Add("Email", "already exists")
	}

	if s.Signup.Password == "" {
		errs.Add("Password", "can't be blank")
	}

	return errs, nil
}

// Signup the struct representing signup attributes
type Signup struct {
	Email    string
	Password string
	Plan     string
	Token    string
}

func SignupsCreate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := SignupRequest{}
	Decode(r, &request)

	errors, err := request.Validate()
	if err != nil {
		return err
	}

	if errors.Any() {
		return errors
	}

	account, err := accountForParams(request.Signup)
	if err != nil {
		return err
	}

	user, err := models.Users.Create(
		request.Signup.Email,
		request.Signup.Password,
		account,
	)
	if err != nil {
		return err
	}

	prefs, err := models.Prefs.Get(user)
	if err != nil {
		logger.Error("getting prefs for a sign up", "err", err, "user", ctx.User.ID, "email", ctx.User.Email)
		return err
	}

	return JSON(w, http.StatusCreated, serializers.NewShowUser(user, prefs))
}

func accountForParams(params Signup) (*models.Account, error) {
	if params.Token != "" {
		invitation, err := models.Invitations.FindByToken(params.Token)
		if err != nil {
			return nil, err
		}

		invitation.Accepted = true
		err = models.Invitations.Update(invitation)
		return &models.Account{ID: invitation.AccountID}, err
	} else {
		account, err := models.Accounts.Create()
		if err != nil {
			return nil, err
		}

		_, err = models.Plans.Create(account, params.Plan)
		if err != nil {
			return nil, err
		}

		return account, nil
	}
}
