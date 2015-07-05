package models

import (
	"time"

	"github.com/erraroo/erraroo/logger"
	"github.com/tuvistavie/securerandom"
)

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
	pr := &PasswordRecover{}
	query := "select * from password_recovers where token = $1 limit 1"
	err := s.Get(pr, query, token)
	if err != nil {
		logger.Error("finding password_recover by token", "token", token, "err", err)
		return nil, err
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
	err := Users.Update(pr.User)
	if err != nil {
		return err
	}

	query := "update password_recovers set used='t' where user_id = $1"
	_, err = s.Exec(query, pr.User.ID)
	if err != nil {
		logger.Error("updating password_recovers", "token", pr.Token, "err", err)
		return err
	}

	return nil
}
