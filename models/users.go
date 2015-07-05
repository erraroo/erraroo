package models

import (
	"database/sql"
	"errors"
)

// UsersStore is the abstraction that needs to be implemented to
// access user data
type UsersStore interface {
	FindByID(int64) (*User, error)
	Exists(email string) bool
	FindByEmail(email string) (*User, error)
	Create(email, password string, account *Account) (*User, error)
	ByAccountID(id int64) ([]*User, error)
	Update(*User) error
}

type usersStore struct {
	*Store
}

// ErrNotFound is returned when a model is not found by the store
var ErrNotFound = errors.New("not found")

func (s usersStore) FindByID(id int64) (*User, error) {
	u := &User{}

	row := s.QueryRow("select id, email, account_id from users where id = $1 limit 1", id)
	err := row.Scan(&u.ID, &u.Email, &u.AccountID)

	if err == sql.ErrNoRows {
		return u, ErrNotFound
	}

	return u, err
}

func (s *usersStore) Insert(user *User) error {
	row := s.QueryRow("insert into users (email, encrypted_password, account_id) values($1,$2,$3) returning id",
		user.Email, user.EncryptedPassword, user.AccountID)
	return row.Scan(&user.ID)
}

func (s *usersStore) Exists(email string) bool {
	exists := false
	row := s.QueryRow("select exists(select 1 from users where lower(email) = lower($1))", email)
	row.Scan(&exists)
	return exists
}

func (s usersStore) FindByEmail(email string) (*User, error) {
	u := &User{}

	row := s.QueryRow("select id, email, encrypted_password, account_id from users where lower(email) = lower($1) limit 1", email)
	err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword, &u.AccountID)

	if err == sql.ErrNoRows {
		return u, ErrNotFound
	}

	return u, err
}

func (s usersStore) Create(email, password string, account *Account) (*User, error) {
	user := NewUser(email, password)
	user.AccountID = account.ID
	return user, s.Insert(user)
}

func (s usersStore) ByAccountID(id int64) ([]*User, error) {
	users := []*User{}
	query := "select * from users where account_id=$1"
	return users, s.Select(&users, query, id)
}

func (s usersStore) Update(user *User) error {
	query := "update users set encrypted_password=$1 where id = $2"
	_, err := s.Exec(query, user.EncryptedPassword, user.ID)
	return err
}
