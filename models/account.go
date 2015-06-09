package models

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
	row := s.QueryRow("insert into accounts default values returning id")
	return account, row.Scan(&account.ID)
}
