package models

type Plan struct {
	AccountID         int64
	RequestsPerMinute int
}

type PlansStore interface {
	FindByToken(string) (*Plan, error)
}

type plansStore struct {
	*Store
}

func (s *plansStore) FindByToken(token string) (*Plan, error) {
	p := &Plan{RequestsPerMinute: 20}
	return p, nil
}
