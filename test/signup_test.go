package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erraroo/erraroo/api/signups"
	"github.com/erraroo/erraroo/models"
	"github.com/stretchr/testify/assert"
)

func TestEmptySignup(t *testing.T) {
	res := httptest.NewRecorder()
	req, res := rr("POST", "/api/v1/signups", struct{}{})

	_app.ServeHTTP(res, req)

	errors := models.ValidationErrors{}
	json.NewDecoder(res.Body).Decode(&errors)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.True(t, errors.Any(), "expected response to have errors")
	assert.Contains(t, errors.Errors["Email"], "can't be blank")
	assert.Contains(t, errors.Errors["Password"], "can't be blank")
}

func TestDuplicateEmailSignup(t *testing.T) {
	signupRequest := signups.SignupRequest{signups.Signup{_user.Email, "password"}}
	req, res := rr("POST", "/api/v1/signups", signupRequest)

	_app.ServeHTTP(res, req)

	errors := models.ValidationErrors{}
	json.NewDecoder(res.Body).Decode(&errors)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.True(t, errors.Any(), "expected response to have errors")
	assert.Contains(t, errors.Errors["Email"], "already exists")
}

func TestValidSignup(t *testing.T) {
	email := uniqEmail()
	signupRequest := signups.SignupRequest{signups.Signup{email, "password"}}
	req, res := rr("POST", "/api/v1/signups", signupRequest)

	_app.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)

	user, err := models.Users.FindByEmail(email)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Nil(t, user.CheckPassword("password"))
}