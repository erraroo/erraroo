package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erraroo/erraroo/api"
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
	signupRequest := api.SignupRequest{api.Signup{_user.Email, "password", "large"}}
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
	signupRequest := api.SignupRequest{api.Signup{email, "password", "large"}}
	req, res := rr("POST", "/api/v1/signups", signupRequest)

	_app.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)

	user, err := models.Users.FindByEmail(email)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Nil(t, user.CheckPassword("password"))

	plan, err := models.Plans.Get(&models.Account{ID: user.AccountID})
	assert.Nil(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, user.AccountID, plan.AccountID)
	assert.Equal(t, 21, plan.DataRetentionInDays)
	assert.Equal(t, 30, plan.RequestsPerMinute)
}
