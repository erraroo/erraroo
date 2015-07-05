package models

import "time"

type PasswordRecover struct {
	Token     string    `db:"token"`
	UserID    int64     `db:"user_id"`
	Used      bool      `db:"used"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	User      *User     `db:"-"`
}
