package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/models"
	"github.com/stretchr/testify/assert"
)

func TestInvalidCreateSession(t *testing.T) {
	errors := models.ValidationErrors{}

	// Empty signin request
	req, res := rr("POST", "/api/v1/sessions", api.SigninRequest{})
	_app.ServeHTTP(res, req)
	json.NewDecoder(res.Body).Decode(&errors)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.True(t, errors.Any(), "expected response to have errors")
	assert.Contains(t, errors.Errors["Signin"], "invalid email or password")

	req, res = rr("POST", "/api/v1/sessions", api.SigninRequest{api.Signin{_user.Email, "INVALID PASSWORD"}})
	_app.ServeHTTP(res, req)
	json.NewDecoder(res.Body).Decode(&errors)
	assert.True(t, errors.Any(), "expected response to have errors")
	assert.Contains(t, errors.Errors["Signin"], "invalid email or password")
}

func TestCreateSession(t *testing.T) {
	req, res := rr("POST", "/api/v1/sessions", api.SigninRequest{api.Signin{_user.Email, "password"}})
	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusCreated, res.Code)

	response := map[string]string{}
	json.NewDecoder(res.Body).Decode(&response)

	req, res = rr("GET", "/api/v1/me", nil)
	req.Header.Add("Authorization", response["token"])

	_app.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}
