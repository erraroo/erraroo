package models

import "errors"

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

	o := s.dbGorm.First(&u, id)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	if o.Error != nil {
		return nil, o.Error
	}

	return u, nil
}

func (s *usersStore) Insert(user *User) error {
	return s.dbGorm.Save(user).Error
}

func (s *usersStore) Exists(email string) bool {
	_, err := s.FindByEmail(email)
	return ErrNotFound != err
}

func (s usersStore) FindByEmail(email string) (*User, error) {
	u := &User{}

	o := s.dbGorm.Where("lower(email) = lower(?)", email).First(&u)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	if o.Error != nil {
		return nil, o.Error
	}

	return u, nil
}

func (s usersStore) Create(email, password string, account *Account) (*User, error) {
	user := NewUser(email, password)
	user.AccountID = account.ID
	return user, s.dbGorm.Save(user).Error
}

func (s usersStore) ByAccountID(id int64) ([]*User, error) {
	users := []*User{}
	return users, s.dbGorm.Where("account_id=?", id).Find(&users).Error
}

func (s usersStore) Update(user *User) error {
	return s.dbGorm.Save(user).Error
}
