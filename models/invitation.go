package models

import "time"

type Invitation struct {
	Token     string `gorm:"primary_key""`
	Address   string
	UserID    int64 `db:"user_id"`
	AccountID int64 `db:"account_id"`
	Accepted  bool
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
