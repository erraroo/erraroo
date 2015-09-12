package models

import (
	"time"

	"github.com/erraroo/erraroo/logger"
)

type RateLimitNotification struct {
	ID        int64
	AccountID int64
	CreatedAt time.Time
}

type RateLimitNotifcationsStore interface {
	Insert(*Account) error
	WasRecentlyNotified(*Account) (bool, error)
}

type rateLimitNotifcationsStore struct{ *Store }

func (store *rateLimitNotifcationsStore) WasRecentlyNotified(account *Account) (bool, error) {
	var count int
	err := store.DB.
		Table("rate_limit_notifications").
		Where("account_id = ?", account.ID).
		Where("created_at >= now_utc() - interval '30 minutes'").Count(&count).Error

	return count > 0, err
}

func (store *rateLimitNotifcationsStore) Insert(account *Account) error {
	n := &RateLimitNotification{
		AccountID: account.ID,
		CreatedAt: time.Now().UTC(),
	}

	query := "insert into rate_limit_notifications (account_id, created_at) values($1, now_utc()) returning id, created_at"
	err := store.QueryRow(query, account.ID).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		logger.Error("creating rate limit notifcation", "account", account.ID, "err", err)
		return err
	}

	return nil
}
