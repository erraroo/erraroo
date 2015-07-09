package models

type PrefsStore interface {
	Get(*User) (*Pref, error)
	Update(*Pref) error
}

type prefsStore struct {
	*Store
}

func (s *prefsStore) Get(u *User) (*Pref, error) {
	pref := new(Pref)
	return pref, s.Where(Pref{UserID: u.ID}).FirstOrCreate(pref).Error
}

func (s *prefsStore) Update(pref *Pref) error {
	return s.Save(pref).Error
}
