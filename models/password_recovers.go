package models

import (
	"time"

	"github.com/erraroo/erraroo/logger"
	"github.com/tuvistavie/securerandom"
)

// PasswordRecoversStore stores password reset tokens for
// users that have forgotten their passwords.
type PasswordRecoversStore interface {
	Create(user *User) (*PasswordRecover, error)
	FindByToken(string) (*PasswordRecover, error)
	Use(*PasswordRecover, string) error
}

type passwordRecoversStore struct{ *Store }

func (s *passwordRecoversStore) Create(user *User) (*PasswordRecover, error) {
	token, err := securerandom.UrlSafeBase64(16, false)
	if err != nil {
		return nil, err
	}

	pr := &PasswordRecover{
		Token:     token,
		UserID:    user.ID,
		User:      user,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	query := "insert into password_recovers (token, user_id, created_at, updated_at) values ($1,$2,$3,$4)"
	_, err = s.Exec(query,
		pr.Token,
		pr.UserID,
		pr.CreatedAt,
		pr.UpdatedAt,
	)

	if err != nil {
		logger.Error("inserting into password_recovers", "err", err)
	}

	return pr, err
}

func (s *passwordRecoversStore) FindByToken(token string) (*PasswordRecover, error) {
	var err error
	pr := &PasswordRecover{}

	o := s.Where("token=?", token).First(&pr)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}

	if o.Error != nil {
		logger.Error("finding password_recover by token", "token", token, "err", o.Error)
		return nil, o.Error
	}

	pr.User, err = Users.FindByID(pr.UserID)
	if err != nil {
		logger.Error("finding user from passswod recover", "token", token, "err", err)
		return nil, err
	}

	return pr, err
}

func (s *passwordRecoversStore) Use(pr *PasswordRecover, pw string) error {
	pr.User.SetPassword(pw)
	if err := Users.Update(pr.User); err != nil {
		return err
	}

	query := "update password_recovers set used='t' where user_id = $1"
	if _, err := s.Exec(query, pr.User.ID); err != nil {
		logger.Error("updating password_recovers", "token", pr.Token, "err", err)
		return err
	}

	return nil
}
