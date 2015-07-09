package models

import "github.com/erraroo/erraroo/logger"

// Account model is the owner of projects
type Account struct {
	ID int64
}

// AccountsStore is the interface to the account records
type AccountsStore interface {
	Create() (*Account, error)
}

type accountsStore struct{ *Store }

func (s *accountsStore) Create() (*Account, error) {
	account := &Account{}
	if err := s.DB.Create(account).Error; err != nil {
		logger.Error("creating account", "err", err)
		return nil, err
	}

	return account, nil
}
