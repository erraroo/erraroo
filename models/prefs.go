package models

import "github.com/erraroo/erraroo/logger"

type PrefsStore interface {
	Get(*User) (*Pref, error)
	Update(*Pref) error
}

type prefsStore struct {
	*Store
}

func (s *prefsStore) Get(u *User) (*Pref, error) {
	query := `
WITH new_row AS (
  INSERT INTO prefs (user_id)
  SELECT $1
  WHERE NOT EXISTS (SELECT * FROM prefs WHERE user_id = $1)
  RETURNING *
)
SELECT * FROM new_row
UNION
SELECT * FROM prefs WHERE user_id = $1;
`
	p := new(Pref)
	err := s.QueryRow(query, u.ID).Scan(
		&p.UserID,
		&p.EmailOnError,
	)

	if err != nil {
		logger.Error("getting prefs", "err", err, "user", u.ID)
	}

	return p, err
}

func (s *prefsStore) Update(pref *Pref) error {
	query := "update prefs set email_on_error=$1 where user_id=$2"
	_, err := s.Exec(query, pref.EmailOnError, pref.UserID)
	return err
}
