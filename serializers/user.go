package serializers

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/erraroo/erraroo/models"
)

type Account struct {
	ID int64
}

type User struct {
	ID        int64
	AccountID int64
	Avatar    string
	Email     string
	PrefID    int64
}

func NewUser(user *models.User) User {
	return User{
		ID:        user.ID,
		AccountID: user.AccountID,
		Avatar:    gravatar(user.Email),
		Email:     user.Email,
		PrefID:    user.ID,
	}
}

type ShowUser struct {
	User     User
	Prefs    []Pref
	Accounts []Account
}

func NewShowUser(user *models.User, pref *models.Pref) ShowUser {
	prefs := []Pref{NewPref(pref)}
	accounts := []Account{Account{user.AccountID}}
	return ShowUser{NewUser(user), prefs, accounts}
}

func gravatar(email string) string {
	hash := md5.New()
	io.WriteString(hash, email)
	return fmt.Sprintf("//www.gravatar.com/avatar/%x?d=identicon", hash.Sum(nil))
}
