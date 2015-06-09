package signups

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
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
}

func Create(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	request := SignupRequest{}
	cx.Decode(r, &request)

	errors, err := request.Validate()
	if err != nil {
		return err
	}

	if errors.Any() {
		return cx.JSON(w, http.StatusBadRequest, errors)
	}

	account, err := models.Accounts.Create()
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

	return cx.JSON(w, http.StatusCreated, serializers.NewShowUser(user))
}
