package models

import "github.com/erraroo/erraroo/logger"

var PlanMap = make(map[string]Plan)

func byName(name string) Plan {
	if plan, ok := PlanMap[name]; ok {
		return plan
	} else {
		return PlanMap["default"]
	}
}

func init() {
	PlanMap["default"] = Plan{RequestsPerMinute: 10, DataRetentionInDays: 7}
	PlanMap["small"] = Plan{RequestsPerMinute: 10, DataRetentionInDays: 7}
	PlanMap["medium"] = Plan{RequestsPerMinute: 20, DataRetentionInDays: 14}
	PlanMap["large"] = Plan{RequestsPerMinute: 30, DataRetentionInDays: 21}
}

type Plan struct {
	AccountID           int64
	DataRetentionInDays int
	RequestsPerMinute   int
}

type PlansStore interface {
	FindByToken(string) (*Plan, error)
	Get(*Account) (*Plan, error)
	Update(*Plan) error

	Create(*Account, string) (*Plan, error)
}

type plansStore struct {
	*Store
}

func (s *plansStore) FindByToken(token string) (*Plan, error) {
	p := &Plan{RequestsPerMinute: 20, DataRetentionInDays: 7}
	return p, nil
}

func (s *plansStore) Create(account *Account, name string) (*Plan, error) {
	model := byName(name)

	plan := &Plan{
		AccountID:           account.ID,
		DataRetentionInDays: model.DataRetentionInDays,
		RequestsPerMinute:   model.RequestsPerMinute,
	}

	query := "insert into plans (account_id, data_retention_in_days, requests_per_minute) values ($1,$2,$3)"
	_, err := s.Exec(query,
		plan.AccountID,
		plan.DataRetentionInDays,
		plan.RequestsPerMinute,
	)

	if err != nil {
		logger.Error("inserting plan", "err", err, "account", account.ID)
	}

	return plan, err
}

func (s *plansStore) Get(account *Account) (*Plan, error) {
	query := `select * from plans where account_id = $1 limit 1;`
	p := new(Plan)
	err := s.QueryRow(query, account.ID).Scan(
		&p.AccountID,
		&p.DataRetentionInDays,
		&p.RequestsPerMinute,
	)

	if err != nil {
		logger.Error("getting plans", "err", err, "account", account.ID)
	}

	return p, err
}

func (s *plansStore) Update(plan *Plan) error {
	query := "update plans set data_retention_in_days=$1, requests_per_minute=$2 where account_id=$3"
	_, err := s.Exec(query, plan.DataRetentionInDays, plan.RequestsPerMinute, plan.AccountID)
	return err
}
