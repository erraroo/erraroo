package models

import (
	"database/sql"
	"sync"
	"time"

	"github.com/erraroo/erraroo/logger"
)

var plans = make(map[string]Plan)

func byName(name string) Plan {
	if plan, ok := plans[name]; ok {
		return plan
	} else {
		return plans["default"]
	}
}

func init() {
	plans["default"] = Plan{RequestsPerMinute: 10, DataRetentionInDays: 7}
	plans["small"] = Plan{RequestsPerMinute: 10, DataRetentionInDays: 7}
	plans["medium"] = Plan{RequestsPerMinute: 20, DataRetentionInDays: 14}
	plans["large"] = Plan{RequestsPerMinute: 30, DataRetentionInDays: 21}
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

	query := "select plans.* from plans join projects on projects.token = $1 and projects.account_id = plans.account_id limit 1;"
	err := s.QueryRow(query, token).Scan(
		&p.AccountID,
		&p.DataRetentionInDays,
		&p.RequestsPerMinute,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		logger.Error("finding plan by token", "token", token, "err", err)
		return nil, err
	}

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
	p := new(Plan)

	query := `select * from plans where account_id = $1 limit 1;`
	err := s.QueryRow(query, account.ID).Scan(
		&p.AccountID,
		&p.DataRetentionInDays,
		&p.RequestsPerMinute,
	)

	if err == sql.ErrNoRows {
		defaultPlan := byName("default")
		logger.Error("using default plan", "account", account.ID)
		return &defaultPlan, nil
	}

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

type planCache struct {
	plans   map[string]*Plan
	expires map[string]time.Time
	lock    *sync.RWMutex
	ttl     time.Duration
}

func NewPlanCache() *planCache {
	return &planCache{
		plans:   map[string]*Plan{},
		expires: map[string]time.Time{},
		lock:    &sync.RWMutex{},
		ttl:     10 * time.Second,
	}
}

func (cache *planCache) FindByToken(token string) (*Plan, error) {
	cache.lock.RLock()
	plan := cache.plans[token]
	expires, cached := cache.expires[token]
	cache.lock.RUnlock()

	if cached && expires.After(time.Now()) {
		return plan, nil
	}

	return cache.put(token)
}

func (cache *planCache) put(token string) (*Plan, error) {
	plan, err := Plans.FindByToken(token)
	if err != nil {
		return nil, err
	}

	cache.lock.Lock()
	cache.plans[token] = plan
	cache.expires[token] = time.Now().Add(cache.ttl)
	cache.lock.Unlock()

	logger.Info("cached plan", "token", token, "plan", plan)
	return plan, nil
}
