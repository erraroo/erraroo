package models

import (
	"errors"
	"reflect"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Encrypter is the password encrypter that should be used
// by default we should use bcrypt
var encrypter PasswordEncrypter

func init() {
	encrypter = &bcryptPasswordEncrypter{}
}

// User model
type User struct {
	ID                int64     `db:"id"`
	Email             string    `db:"email"`
	EncryptedPassword []byte    `db:"encrypted_password"`
	AccountID         int64     `db:"account_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// CheckPassword ensures that the user's password matches
func (user User) CheckPassword(password string) error {
	return encrypter.Check(user.EncryptedPassword, password)
}

// SetPassword takes a clear text password and encrypts it
func (user *User) SetPassword(password string) (err error) {
	user.EncryptedPassword, err = encrypter.Encrypt(password)
	return err
}

func (user *User) CanSee(model interface{}) bool {
	r := reflect.ValueOf(model)
	f := reflect.Indirect(r).FieldByName("AccountID")
	return user.AccountID == f.Int()
}

// NewUser returns a new user with an encrypted password
func NewUser(email, password string) *User {
	user := &User{Email: email}
	user.SetPassword(password)
	return user
}

// PasswordEncrypter encrypts passwords
type PasswordEncrypter interface {
	Encrypt(string) ([]byte, error)
	Check(hashed []byte, password string) error
}

type bcryptPasswordEncrypter struct{}

func (b *bcryptPasswordEncrypter) Encrypt(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (b *bcryptPasswordEncrypter) Check(hashed []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashed, []byte(password))
}

type dummyPasswordEncrypter struct{}

func (b *dummyPasswordEncrypter) Encrypt(password string) ([]byte, error) {
	return []byte(password), nil
}

func (b *dummyPasswordEncrypter) Check(hashed []byte, password string) error {
	if password == string(hashed) {
		return nil
	}

	return errors.New("does not match")
}
