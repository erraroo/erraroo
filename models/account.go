package models

import "github.com/erraroo/erraroo/logger"

// Account model is the owner of projects
type Account struct {
	ID int64
}

func CreateAccount() (*Account, error) {
	account := &Account{}
	if err := store.Create(account).Error; err != nil {
		logger.Error("creating account", "err", err)
		return nil, err
	}

	return account, nil
}
