package test

import (
	"testing"

	"github.com/erraroo/erraroo/api"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvitation(t *testing.T) {
	request := api.CreateInvitationRequest{}
	request.Invitation.To = "ben@nerdyworm.com"

	req, res := rr("POST", "/api/v1/invitations", request)
	req.Header.Add("Authorization", _token)

	_app.ServeHTTP(res, req)
	assert.Equal(t, 201, res.Code)
}
