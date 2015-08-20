package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
)

// SigninRequest is an incoming sign in request
type SigninRequest struct {
	Signin Signin
}

// Signin the params for a sign in
type Signin struct {
	Email    string
	Password string
}

type Success struct {
	Token  string `json:"token"`
	UserID string `json:"userID"`
}

func SessionsCreate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	errors := models.NewValidationErrors()
	errors.Add("Signin", "invalid email or password")

	request := SigninRequest{}
	Decode(r, &request)

	if request.Signin.Email == "" || request.Signin.Password == "" {
		return errors
	}

	user, err := models.Users.FindByEmail(request.Signin.Email)
	if err == models.ErrNotFound {
		return errors
	} else if err != nil {
		return err
	}

	err = user.CheckPassword(request.Signin.Password)
	if err != nil {
		return errors
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["user_id"] = user.ID
	token.Claims["expires"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(config.TokenSigningKey)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusCreated, Success{tokenString, fmt.Sprintf("%d", user.ID)})
}

func SessionsDestroy(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	return nil
}
