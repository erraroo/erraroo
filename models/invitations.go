package models

import (
	"time"

	"github.com/erraroo/erraroo/logger"
	"github.com/tuvistavie/securerandom"
)

type InvitationsStore interface {
	ListForUser(*User) ([]*Invitation, error)
	Create(address string, user *User) (*Invitation, error)
	FindByToken(string) (*Invitation, error)
	Update(*Invitation) error
}

type invitationsStore struct{ *Store }

func (s *invitationsStore) ListForUser(user *User) ([]*Invitation, error) {
	invitations := []*Invitation{}

	query := "select * from invitations where account_id = $1 order by created_at desc, updated_at desc"
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

	query := "insert into invitations (token, user_id, account_id, address, created_at, updated_at) values ($1,$2,$3,$4,$5,$6)"
	_, err = s.Exec(query,
		invitation.Token,
		invitation.UserID,
		invitation.AccountID,
		invitation.Address,
		invitation.CreatedAt,
		invitation.UpdatedAt,
	)

	if err != nil {
		logger.Error("inserting into invitations", "err", err)
	}

	return invitation, err
}

func (s *invitationsStore) FindByToken(token string) (*Invitation, error) {
	invitation := &Invitation{}
	query := "select * from invitations where token = $1 limit 1"
	err := s.Get(invitation, query, token)
	if err != nil {
		logger.Error("finding invitation by token", "token", token, "err", err)
	}

	return invitation, err
}

func (s *invitationsStore) Update(i *Invitation) error {
	i.UpdatedAt = time.Now().UTC()
	query := "update invitations set accepted=$1, updated_at=$2 where token = $3"

	_, err := s.Exec(query, i.Accepted, i.UpdatedAt, i.Token)
	if err != nil {
		logger.Error("updating invitation", "token", i.Token, "err", err)
		return err
	}

	return nil
}
