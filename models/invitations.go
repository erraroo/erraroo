package models

import (
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
	return invitations, s.Where("account_id=?", user.AccountID).
		Order("created_at desc").Find(&invitations).Error
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
	}

	return invitation, s.DB.Create(invitation).Error
}

func (s *invitationsStore) FindByToken(token string) (*Invitation, error) {
	invitation := &Invitation{}

	o := s.Where("token=?", token).First(&invitation)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	if o.Error != nil {
		logger.Error("finding invitation by token", "token", token, "err", o.Error)
	}

	return invitation, nil
}

func (s *invitationsStore) Update(i *Invitation) error {
	if err := s.Save(i).Error; err != nil {
		logger.Error("updating invitation", "token", i.Token, "err", err)
		return err
	}

	return nil
}
