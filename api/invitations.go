package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
)

type CreateInvitationRequest struct {
	Invitation InvitationParams
}

type InvitationParams struct {
	To string
}

func InvitationsCreate(w http.ResponseWriter, r *http.Request, c *cx.Context) error {
	return JSON(w, http.StatusCreated, nil)
}
