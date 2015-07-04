package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/gorilla/mux"
)

type CreateInvitationRequest struct {
	Invitation InvitationParams
}

type InvitationParams struct {
	Address string
}

func InvitationsCreate(w http.ResponseWriter, r *http.Request, c *cx.Context) error {
	request := CreateInvitationRequest{}
	Decode(r, &request)

	errors := models.NewValidationErrors()
	if request.Invitation.Address == "" {
		errors.Add("Address", "can not be blank")
	}

	if errors.Any() {
		return errors
	}

	invitation, err := models.Invitations.Create(request.Invitation.Address, c.User)
	if err != nil {
		logger.Error("creating invitation", "err", err)
	}

	return JSON(w, http.StatusCreated, serializers.NewShowInvitation(invitation))
}

func InvitationsShow(w http.ResponseWriter, r *http.Request, c *cx.Context) error {
	token := mux.Vars(r)["token"]
	logger.Info("show token", "token", token)

	invite, err := models.Invitations.FindByToken(token)
	if err != nil {
		logger.Error("finding token", "err", err)
		return nil
	}

	return JSON(w, http.StatusOK, serializers.NewShowInvitation(invite))
}
