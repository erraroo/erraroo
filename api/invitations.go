package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/erraroo/erraroo/usecases"
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

	invitation, err := usecases.InviteByEmail(c.User, request.Invitation.Address)
	if err != nil {
		logger.Error("usecases.InviteByEmail", "err", err)
		return err
	}

	return JSON(w, http.StatusCreated, serializers.NewShowInvitation(invitation))
}

func InvitationsShow(w http.ResponseWriter, r *http.Request, c *cx.Context) error {
	token := mux.Vars(r)["token"]
	invite, err := models.Invitations.FindByToken(token)
	if err != nil {
		logger.Error("finding token", "err", err, "token", token)
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowInvitation(invite))
}

func InvitationsIndex(w http.ResponseWriter, r *http.Request, c *cx.Context) error {
	invitations, err := models.Invitations.ListForUser(c.User)
	if err != nil {
		logger.Error("listing invitations", "err", err)
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewInvitations(invitations))
}

func InvitationsDelete(w http.ResponseWriter, r *http.Request, c *cx.Context) error {
	token := mux.Vars(r)["token"]
	invite, err := models.Invitations.FindByToken(token)
	if !c.User.CanSee(invite) {
		return models.ErrNotFound
	}

	err = models.Invitations.Delete(invite)
	if err != nil {
		logger.Error("deleting an invitation", "err", err)
		return err
	}

	return JSON(w, http.StatusNoContent, nil)
}
