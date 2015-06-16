package serializers

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/erraroo/erraroo/models"
)

type User struct {
	ID     int64
	Email  string
	Avatar string
	PrefID int64
}

func NewUser(user *models.User) User {
	return User{
		ID:     user.ID,
		Email:  user.Email,
		Avatar: gravatar(user.Email),
		PrefID: user.ID,
	}
}

type ShowUser struct {
	User  User
	Prefs []Pref
}

func NewShowUser(user *models.User, pref *models.Pref) ShowUser {
	prefs := []Pref{NewPref(pref)}
	return ShowUser{NewUser(user), prefs}
}

func gravatar(email string) string {
	hash := md5.New()
	io.WriteString(hash, email)
	return fmt.Sprintf("http://www.gravatar.com/avatar/%x?d=identicon", hash.Sum(nil))
}
