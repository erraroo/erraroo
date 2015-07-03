package models

import (
	"time"

	"github.com/erraroo/erraroo/logger"
	"github.com/tuvistavie/securerandom"
)

type InvitationsStore interface {
	ListForUser(*User) ([]*Invitation, error)
	Create(address string, user *User) (*Invitation, error)
}

type invitationsStore struct{ *Store }

func (s *invitationsStore) ListForUser(user *User) ([]*Invitation, error) {
	invitations := []*Invitation{}

	query := "select * from invitations where account_id = $1"
	err := s.Select(&invitations, query, user.AccountID)
	if err != nil {
		logger.Error("selecting from invitations", "err", err)
	}

	return invitations, err
}

func (s *invitationsStore) Create(to string, user *User) (*Invitation, error) {
	token, err := securerandom.UrlSafeBase64(16, false)
	if err != nil {
		return nil, err
	}

	invitation := &Invitation{
		Token:     token,
		UserID:    user.ID,
		AccountID: user.AccountID,
		Address:   to,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	query := "insert into invitations (token, user_id, account_id, address) values ($1,$2,$3,$4)"
	_, err = s.Exec(query,
		invitation.Token,
		invitation.UserID,
		invitation.AccountID,
		invitation.Address,
	)

	if err != nil {
		logger.Error("inserting into invitations", "err", err)
	}

	return invitation, err
}
